package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/ethanhosier/worker-node/worker_manager"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Add command line flags
	workerType := flag.String("worker", "", "Type of worker to run (scraper or rag)")
	concurrency := flag.Int("concurrency", 0, "Number of concurrent workers (only used for scraper)")
	flag.Parse()

	if *workerType == "" {
		log.Fatal("Please specify a worker type using -worker flag (scraper or rag)")
	}

	if *workerType == "scraper" && *concurrency == 0 {
		log.Fatal("Please specify a concurrency using -concurrency flag (scraper)")
	}

	var (
		coordinatorClient = coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
	)

	switch *workerType {
	case "scraper":
		_, errCh := startScraperWorkerManager(coordinatorClient, *concurrency)
		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	case "rag":
		_, errCh := startRagWorkerManager(coordinatorClient)
		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	default:
		log.Fatalf("Unknown worker type: %s", *workerType)
	}

	select {}
}

func startScraperWorkerManager(coordinatorClient coordinator_client.CoordinatorClient, concurrency int) (chan<- bool, <-chan error) {
	var (
		scraperClient = scraper.NewHttpScraper()

		scraperWorkerManager = worker_manager.NewScraperWorkerManager(context.TODO(), coordinatorClient, scraperClient, concurrency)
	)

	return scraperWorkerManager.Start()
}

func startRagWorkerManager(coordinatorClient coordinator_client.CoordinatorClient) (chan<- bool, <-chan error) {
	var (
		ragClient        = ragger.NewRagClient("./model", "./libonnxruntime.so.1.20.1")
		store            = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
		ragWorkerManager = worker_manager.NewRagWorkerManager(context.TODO(), coordinatorClient, ragClient, store, 1)
	)

	return ragWorkerManager.Start()
}
