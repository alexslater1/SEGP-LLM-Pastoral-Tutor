package worker_manager

import (
	"context"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
)

type WorkerConfigType string

const (
	WorkerConfigTypeScraper WorkerConfigType = "scraper"
	WorkerConfigTypeRag     WorkerConfigType = "rag"
)

type WorkerConfig struct {
	Type              WorkerConfigType
	ctx               context.Context
	coordinatorClient coordinator_client.CoordinatorClient
	numWorkers        int

	scraper scraper.Scraper

	ragger ragger.Ragger
	store  storage.Storage
}

func NewScraperWorkerManager(ctx context.Context, coordinatorClient coordinator_client.CoordinatorClient, scraper scraper.Scraper, numWorkers int) *WorkerManager {
	workerConfig := &WorkerConfig{
		Type:              WorkerConfigTypeScraper,
		ctx:               ctx,
		coordinatorClient: coordinatorClient,
		scraper:           scraper,
		numWorkers:        numWorkers,
	}

	return newWorkerManager(workerConfig)
}

func NewRagWorkerManager(ctx context.Context, coordinatorClient coordinator_client.CoordinatorClient, ragger ragger.Ragger, store storage.Storage, numWorkers int) *WorkerManager {
	workerConfig := &WorkerConfig{
		Type:              WorkerConfigTypeRag,
		ctx:               ctx,
		coordinatorClient: coordinatorClient,
		ragger:            ragger,
		store:             store,
		numWorkers:        numWorkers,
	}

	return newWorkerManager(workerConfig)
}
