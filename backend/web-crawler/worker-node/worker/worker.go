package worker

import (
	"context"

	coordinator_client "github.com/ethanhosier/worker-node/coordinator_client"
)

type WorkerType string

const (
	WorkerTypeScraper WorkerType = "scraper"
	WorkerTypeRag     WorkerType = "rag"
)

type Worker interface {
	Execute(ctx context.Context, task *coordinator_client.Task) error
	Cleanup(ctx context.Context, task *coordinator_client.Task) error
	WorkerType() WorkerType
	Id() string
}
