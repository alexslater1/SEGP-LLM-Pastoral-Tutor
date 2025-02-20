package jobs

import (
	"log"
	"time"
)

type JobManager struct {
	jobs []Job
}

func NewJobManager(jobs []Job) *JobManager {
	return &JobManager{jobs: jobs}
}

func (j *JobManager) Start() <-chan error {
	errCh := make(chan error)
	tickers := make([]*time.Ticker, len(j.jobs))
	for i, job := range j.jobs {
		tickers[i] = time.NewTicker(job.Interval())

		// Start a goroutine for each job
		go func(job Job, ticker *time.Ticker) {
			for range ticker.C {
				log.Printf("Running job %s", job.Name())
				err := job.Run()
				if err != nil {
					log.Printf("Job %s failed: %v. Continuing", job.Name(), err)
					errCh <- err
				}
				log.Printf("Job %s completed", job.Name())
			}
		}(job, tickers[i])
	}

	return errCh
}
