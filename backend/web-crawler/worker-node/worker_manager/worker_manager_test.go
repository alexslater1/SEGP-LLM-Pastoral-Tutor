package worker_manager

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/ethanhosier/worker-node/storage"
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

func TestWorkerManagerScraper(t *testing.T) {
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

func TestWorkerManagerRag(t *testing.T) {
	// given
	var (
		ragClient         = ragger.NewMockRagClient()
		coordinatorClient = coordinator_client.NewMockCoordinatorClient()
		store             = storage.NewMemoryStorage()

		workerManager = NewRagWorkerManager(context.TODO(), coordinatorClient, ragClient, store, 1)

		markdown       = "# Hello, World!"
		text           = "Hello, World!"
		mockTask1, err = coordinator_client.NewTask("1", "CREATED_BY", worker.RagWorkerParams{
			Markdown:  markdown,
			Url:       "https://example.com",
			InnerText: text,
		})

		chunks1     = []string{"Hello, World!1", "Hello, World!2"}
		embeddings1 = [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}, {7.0, 8.0, 9.0}}
	)
	if err != nil {
		t.Fatalf("Error creating mock task: %v", err)
	}

	ragClient.SetChunksFor(text, chunks1)

	ragClient.SetContactsFor(markdown, []ragger.Contact{
		{
			Value:   "johndoe@example.com",
			Context: markdown,
			Type:    "email",
		},
	})

	newSlice := make([]string, len(chunks1)+1)
	copy(newSlice, chunks1)
	newSlice[len(chunks1)] = markdown

	ragClient.SetEmbeddingsForAll(newSlice, embeddings1)

	err = coordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicRag, mockTask1)
	if err != nil {
		t.Fatalf("Error creating mock task: %v", err)
	}

	// when
	doneChan, errChan := workerManager.Start()

	time.Sleep(10 * time.Second)

	doneChan <- true

	// then
	select {
	case err := <-errChan:
		t.Fatalf("Error: %v", err)
	default:
	}

	_, err = coordinatorClient.GetTask(context.TODO(), 1*time.Second, coordinator_client.CoordinatorClientTaskTopicRag)
	if err != coordinator_client.ErrNoTasksToComplete {
		t.Fatal("There should be no tasks to complete error")
	}

	storedRagSources, err := storage.GetAll[storage.RagSource](store, nil)
	if err != nil {
		t.Fatalf("Error getting stored rag source: %v", err)
	}

	assert.Equal(t, len(storedRagSources), 1)
	assert.Equal(t, storedRagSources[0].URL, "https://example.com")

	contacts, err := storage.GetAll[storage.RagContact](store, nil)
	if err != nil {
		t.Fatalf("Error getting contacts: %v", err)
	}

	assert.Equal(t, len(contacts), 1)
	assert.Equal(t, contacts[0].Context, markdown)
	assert.Equal(t, contacts[0].Contact, "johndoe@example.com")
	assert.Equal(t, contacts[0].ContactType, "email")
	assert.Equal(t, contacts[0].RagSourceId, storedRagSources[0].ID)
	assert.Equal(t, contacts[0].Embedding, []float32{7.0, 8.0, 9.0})

	rags, err := storage.GetAll[storage.RagChunk](store, nil)
	if err != nil {
		t.Fatalf("Error getting rags: %v", err)
	}

	assert.Equal(t, len(rags), 2)
	assert.Equal(t, rags[0].Text, chunks1[0])
	assert.Equal(t, rags[1].Text, chunks1[1])
	assert.Equal(t, rags[0].Embedding, embeddings1[0])
	assert.Equal(t, rags[1].Embedding, embeddings1[1])
	assert.Equal(t, rags[0].RagSourceId, storedRagSources[0].ID)
	assert.Equal(t, rags[1].RagSourceId, storedRagSources[0].ID)
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
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test because CICD is true")
	}

	scraper := scraper.NewHttpScraper()
	redisClient := coordinator_client.NewRedisCoordinatorClient(context.TODO(), "localhost:6379", "", 0)
	scraperWorkerManager := NewScraperWorkerManager(context.TODO(), redisClient, scraper, 1)

	workers := scraperWorkerManager.createWorkers()

	if len(workers) != 1 {
		t.Fatalf("Expected 1 worker, got %d", len(workers))
	}
}

func TestRedisWorkerManager2(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test because CICD is true")
	}

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
