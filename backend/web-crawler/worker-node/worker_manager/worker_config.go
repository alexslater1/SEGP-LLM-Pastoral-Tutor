package worker_manager

import (
	"context"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/scraper"
)

type WorkerConfigType string

const (
	WorkerConfigTypeScraper WorkerConfigType = "scraper"
)

type WorkerConfig struct {
	Type              WorkerConfigType
	ctx               context.Context
	coordinatorClient coordinator_client.CoordinatorClient
	numWorkers        int
	taskQueue         string
	processingQueue   string

	scraper scraper.Scraper
}

func NewScraperWorkerManager(ctx context.Context, coordinatorClient coordinator_client.CoordinatorClient, scraper scraper.Scraper, numWorkers int) *WorkerManager {
	workerConfig := &WorkerConfig{
		Type:              WorkerConfigTypeScraper,
		ctx:               ctx,
		coordinatorClient: coordinatorClient,
		scraper:           scraper,
		numWorkers:        numWorkers,

		processingQueue: "processing_queue",
		taskQueue:       "url_queue",
	}

	return newWorkerManager(workerConfig)
}
