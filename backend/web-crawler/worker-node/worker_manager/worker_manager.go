package worker_manager

import (
	"log"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/worker"
)

const (
	getTaskTimeout = 3 * time.Second
)

type WorkerManager struct {
	config *WorkerConfig
}

func newWorkerManager(config *WorkerConfig) *WorkerManager {
	return &WorkerManager{
		config: config,
	}
}

func (w *WorkerManager) Start() (chan<- bool, <-chan error) {
	workers := w.createWorkers()

	errChan := make(chan error)
	taskChan := make(chan *coordinator_client.Task)
	doneChan := make(chan bool)

	go func() {
		err := w.TaskLoop(taskChan, doneChan)
		if err != nil {
			errChan <- err
		}

		close(taskChan)
	}()

	for _, worker := range workers {
		go w.workerLoop(worker, taskChan, errChan)
	}

	return doneChan, errChan
}

func (w *WorkerManager) TaskLoop(taskChan chan<- *coordinator_client.Task, doneCh <-chan bool) error {
	for {
		// Check if we should stop
		select {
		case <-doneCh:
			return nil
		default:
			// Continue with task fetching
		}

		// Try to get a task
		task, err := w.config.coordinatorClient.GetTaskAndSetProcessing(w.config.ctx, getTaskTimeout, coordinator_client.CoordinatorClientTaskTopicUrls)

		if err == coordinator_client.ErrNoTasksToComplete {
			log.Println("No tasks to complete, waiting for new tasks...")
			continue
		}

		if err != nil {
			return err
		}

		log.Printf("Task Found: %s", task.ID)
		taskChan <- task
	}
}

func (w *WorkerManager) workerLoop(worker worker.Worker, taskChan <-chan *coordinator_client.Task, errorChan chan<- error) {
	log.Printf("Worker %s starting", worker.Id())

	for task := range taskChan {
		log.Printf("Worker %s executing task %s", worker.Id(), task.ID)
		err := worker.Execute(w.config.ctx, task)
		if err != nil {
			errorChan <- err
			return
		}

		log.Printf("Worker %s cleaning up task %s", worker.Id(), task.ID)
		err = worker.Cleanup(w.config.ctx, task)
		if err != nil {
			errorChan <- err
			return
		}
	}
}

func (w *WorkerManager) createWorkers() []worker.Worker {
	workers := make([]worker.Worker, w.config.numWorkers)

	for i := 0; i < w.config.numWorkers; i++ {
		switch w.config.Type {
		case WorkerConfigTypeScraper:
			workers[i] = worker.NewScraperWorker(w.config.scraper, w.config.coordinatorClient)
		}
	}

	return workers
}
