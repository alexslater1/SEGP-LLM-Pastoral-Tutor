package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/ethanhosier/worker-node/worker_manager"
	"github.com/joho/godotenv"
)

var (
	modelPath     = filepath.Join("model", "model.onnx")
	libraryPath   = filepath.Join("libonnxruntime.so.1.20.1")
	tokenizerPath = filepath.Join("model", "tokenizer.json")
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Add command line flags
	workerType := flag.String("worker", "", "Type of worker to run (scraper or rag)")
	concurrency := flag.Int("concurrency", 0, "Number of concurrent workers (only used for scraper)")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis server address")
	redisPassword := flag.String("redis-password", "", "Redis password")
	redisDB := flag.Int("redis-db", 0, "Redis database number")
	flag.Parse()

	if *workerType == "" {
		log.Fatal("Please specify a worker type using -worker flag (scraper or rag)")
	}

	if *workerType == "scraper" && *concurrency == 0 {
		log.Fatal("Please specify a concurrency using -concurrency flag (scraper)")
	}

	var (
		coordinatorClient = coordinator_client.NewRedisCoordinatorClient(
			context.TODO(),
			*redisAddr,
			*redisPassword,
			*redisDB,
		)
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
		ragClient        = ragger.NewRAGClient(modelPath, libraryPath, tokenizerPath)
		store            = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
		ragWorkerManager = worker_manager.NewRagWorkerManager(context.TODO(), coordinatorClient, ragClient, store, 1)
	)

	return ragWorkerManager.Start()
}
