package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/ethanhosier/worker-node/utils"
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

	redisAddr := utils.Required(os.Getenv("REDIS_ADDR"), "REDIS_ADDR")
	redisPassword := utils.Required(os.Getenv("REDIS_PASSWORD"), "REDIS_PASSWORD")
	redisDB := utils.RequiredInt(os.Getenv("REDIS_DB"), "REDIS_DB")

	log.Printf("Redis address: %s, password: %s, db: %d", redisAddr, redisPassword, redisDB)

	workerType := utils.Required(os.Getenv("WORKER_TYPE"), "WORKER_TYPE")

	var (
		coordinatorClient = coordinator_client.NewRedisCoordinatorClient(
			context.TODO(),
			redisAddr,
			redisPassword,
			redisDB,
		)
	)

	switch workerType {
	case "scraper":
		concurrency := utils.RequiredInt(os.Getenv("CONCURRENCY"), "CONCURRENCY")

		_, errCh := startScraperWorkerManager(coordinatorClient, concurrency)
		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	case "rag":
		_, errCh := startRagWorkerManager(coordinatorClient)
		for err := range errCh {
			log.Fatalf("Error: %v", err)
		}
	default:
		log.Fatalf("Unknown worker type: %s", workerType)
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
