package worker_manager

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/worker"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestWorkerManager(t *testing.T) {
	var (
		scraper              = scraper.NewMockScraper()
		coordinatorClient    = coordinator_client.NewMockCoordinatorClient()
		scraperWorkerManager = NewScraperWorkerManager(context.TODO(), coordinatorClient, scraper, 1)

		mockTask1, err = coordinator_client.NewTask("1", "CREATED_BY", worker.ScraperWorkerParams{
			Url: "https://example.com",
		})
	)
	if err != nil {
		t.Fatalf("Error creating mock task: %v", err)
	}

	scraper.SetHtmlContent("https://example.com", "<html><body><main><h1>Hello, World!</h1></main></body></html>")

	coordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicUrls, mockTask1)

	panic(scraperWorkerManager.Start())
}

func TestWorkerManagerRedis(t *testing.T) {
	var (
		scraper              = scraper.NewHttpScraper()
		redisClient          = coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
		scraperWorkerManager = NewScraperWorkerManager(context.TODO(), redisClient, scraper, 1)
	)

	panic(scraperWorkerManager.Start())
}
