package worker_manager

import (
	"fmt"
	"log"
	"strings"
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
		task, err := w.config.coordinatorClient.GetTaskAndSetProcessing(w.config.ctx, getTaskTimeout, topicForWorkerConfigType(w.config.Type))

		if err == coordinator_client.ErrNoTasksToComplete {
			log.Printf("No %s tasks to complete, waiting for new tasks...", strings.ToUpper(string(w.config.Type)))
			continue
		}

		if err != nil {
			return err
		}

		log.Printf("%s Task Found: %s", strings.ToUpper(string(w.config.Type)), task.ID)
		taskChan <- task
	}
}

func (w *WorkerManager) workerLoop(worker worker.Worker, taskChan <-chan *coordinator_client.Task, errorChan chan<- error) {
	log.Printf("%s Worker %s starting", strings.ToUpper(string(w.config.Type)), worker.Id())

	for task := range taskChan {
		log.Printf("%s Worker %s executing task %s", strings.ToUpper(string(w.config.Type)), worker.Id(), task.ID)
		err := worker.Execute(w.config.ctx, task)
		if err != nil {
			log.Printf("%s Worker %s failed to execute task %s: %v. Will store error and continue.",
				strings.ToUpper(string(w.config.Type)), worker.Id(), task.ID, err)

			err = w.config.coordinatorClient.StoreError(w.config.ctx, topicForWorkerConfigType(w.config.Type), task, err)
			if err != nil {
				log.Printf("Failed to store error for task %s: %v", task.ID, err)
				errorChan <- err
				return
			}
		}

		log.Printf("%s Worker %s cleaning up task %s", strings.ToUpper(string(w.config.Type)), worker.Id(), task.ID)
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
		case WorkerConfigTypeRag:
			workers[i] = worker.NewRagWorker(w.config.ragger, w.config.coordinatorClient, w.config.store)
		}
	}

	return workers
}

func topicForWorkerConfigType(workerType WorkerConfigType) coordinator_client.CoordinatorClientTaskTopic {
	switch workerType {
	case WorkerConfigTypeScraper:
		return coordinator_client.CoordinatorClientTaskTopicUrls
	case WorkerConfigTypeRag:
		return coordinator_client.CoordinatorClientTaskTopicRag
	}

	panic(fmt.Sprintf("Unknown worker type: %s", workerType))
}
