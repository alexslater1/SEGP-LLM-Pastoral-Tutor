package main

import (
	"context"
	"log"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/ethanhosier/worker-node/worker_manager"
)

func main() {
	var (
		coordinatorClient = coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
	)

	go func() {
		_, errCh := startScraperWorkerManager(coordinatorClient)

		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	}()

	go func() {
		_, errCh := startRagWorkerManager(coordinatorClient)

		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	}()

	select {}
}

func startScraperWorkerManager(coordinatorClient coordinator_client.CoordinatorClient) (chan<- bool, <-chan error) {
	var (
		scraperClient = scraper.NewHttpScraper()

		scraperWorkerManager = worker_manager.NewScraperWorkerManager(context.TODO(), coordinatorClient, scraperClient, 1)
	)

	return scraperWorkerManager.Start()
}

func startRagWorkerManager(coordinatorClient coordinator_client.CoordinatorClient) (chan<- bool, <-chan error) {
	var (
		ragClient        = ragger.NewRagClient("./model", "./libonnxruntime.so.1.20.1")
		store            = storage.NewMemoryStorage()
		ragWorkerManager = worker_manager.NewRagWorkerManager(context.TODO(), coordinatorClient, ragClient, store, 1)
	)

	return ragWorkerManager.Start()
}
