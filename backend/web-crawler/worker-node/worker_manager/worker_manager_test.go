package worker_manager

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/worker"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

	err = coordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicUrls, mockTask1)
	if err != nil {
		t.Fatalf("Error creating mock task: %v", err)
	}

	doneChan, errChan := scraperWorkerManager.Start()

	time.Sleep(10 * time.Second)

	doneChan <- true

	select {
	case err := <-errChan:
		t.Fatalf("Error: %v", err)
	default:
	}

	ragTask, err := coordinatorClient.GetTask(context.TODO(), 1*time.Second, coordinator_client.CoordinatorClientTaskTopicRag)
	if err != nil {
		t.Fatalf("Error getting rag task: %v", err)
	}

	parsedParams, err := coordinator_client.CastParams[worker.RagWorkerParams](ragTask.Params)
	if err != nil {
		t.Fatalf("Error parsing rag task params: %v", err)
	}

	assert.Equal(t, parsedParams.Url, "https://example.com")
	assert.Equal(t, parsedParams.Markdown, "# Hello, World!")
}

func TestWorkerManagerRedis(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test because CICD is true")
	}

	var (
		scraper              = scraper.NewHttpScraper()
		redisClient          = coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
		scraperWorkerManager = NewScraperWorkerManager(context.TODO(), redisClient, scraper, 1)
	)

	doneChan, errChan := scraperWorkerManager.Start()

	time.Sleep(5 * time.Second)

	doneChan <- true

	select {
	case err := <-errChan:
		t.Fatalf("Error:  %v", err)
	default:
	}
}

func TestCreateWorkers(t *testing.T) {
	scraper := scraper.NewHttpScraper()
	redisClient := coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
	scraperWorkerManager := NewScraperWorkerManager(context.TODO(), redisClient, scraper, 1)

	workers := scraperWorkerManager.createWorkers()

	if len(workers) != 1 {
		t.Fatalf("Expected 1 worker, got %d", len(workers))
	}
}

func TestRedisWorkerManager2(t *testing.T) {
	var (
		scraper                = scraper.NewHttpScraper()
		redisCoordinatorClient = coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
		scraperWorkerManager   = NewScraperWorkerManager(context.TODO(), redisCoordinatorClient, scraper, 1)

		mockTask1, err = coordinator_client.NewTask("1", "CREATED_BY", worker.ScraperWorkerParams{
			Url: "https://www.ethanhosier.com/",
		})
	)
	if err != nil {
		t.Fatalf("Error creating mock task: %v", err)
	}

	err = redisCoordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicUrls, mockTask1)
	if err != nil {
		t.Fatalf("Error creating mock  task: %v", err)
	}

	doneChan, errChan := scraperWorkerManager.Start()

	time.Sleep(10 * time.Second)

	doneChan <- true

	select {
	case err := <-errChan:
		t.Fatalf("Error: %v", err)
	default:
	}

	ragTask, err := redisCoordinatorClient.GetTask(context.TODO(), 1*time.Second, coordinator_client.CoordinatorClientTaskTopicRag)
	if err != nil {
		t.Fatalf("Error getting rag task: %v", err)
	}

	parsedParams, err := coordinator_client.CastParams[worker.RagWorkerParams](ragTask.Params)
	if err != nil {
		t.Fatalf("Error parsing rag task params: %v", err)
	}

	assert.Equal(t, parsedParams.Url, "https://www.ethanhosier.com/")
	assert.NotNil(t, parsedParams.Markdown)

	t.Logf("Rag task: %v", ragTask)
}
