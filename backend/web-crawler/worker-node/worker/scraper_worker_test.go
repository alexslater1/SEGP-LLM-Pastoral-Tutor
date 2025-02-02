package worker

import (
	"context"
	"testing"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/scraper"
	"github.com/stretchr/testify/assert"
)

func TestScraperWorkerMdAndTextFromUrl(t *testing.T) {
	scraper := scraper.NewMockScraper()
	scraper.SetHtmlContent("https://example.com", "<html><body><main>Hello, world!</main></body></html>")

	worker := NewScraperWorker(scraper, nil)
	md, text, err := worker.mdAndTextFromUrl("https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, md, "Hello, world!")
	assert.Equal(t, text, "Hello, world!")
}

func TestScraperWorkerMdFromUrlContent(t *testing.T) {
	scraper := scraper.NewMockScraper()
	worker := NewScraperWorker(scraper, nil)

	// Test nested content
	scraper.SetHtmlContent("https://nested.com", `
		<html><body><main>
			<h1>Title</h1>
			<p>Paragraph 1</p>
			<div>
				<p>Nested paragraph</p>
			</div>
		</main></body></html>
	`)
	md, text, err := worker.mdAndTextFromUrl("https://nested.com")
	assert.NoError(t, err)
	assert.Contains(t, md, "Title")
	assert.Contains(t, md, "Paragraph 1")
	assert.Contains(t, md, "Nested paragraph")
	assert.Contains(t, text, "Title")
	assert.Contains(t, text, "Paragraph 1")
	assert.Contains(t, text, "Nested paragraph")

	// Test multiple main tags
	scraper.SetHtmlContent("https://multiplemain.com", `
		<html><body>
			<main>First main</main>
		</body></html>
	`)
	md, text, err = worker.mdAndTextFromUrl("https://multiplemain.com")
	assert.NoError(t, err)
	assert.Contains(t, md, "First main")
	assert.Contains(t, text, "First main")
}

func TestScraperWorkerId(t *testing.T) {
	scraper := scraper.NewMockScraper()
	worker := NewScraperWorker(scraper, nil)
	assert.NotEmpty(t, worker.Id())
}

func TestScraperWorkerType(t *testing.T) {
	scraper := scraper.NewMockScraper()
	worker := NewScraperWorker(scraper, nil)
	assert.Equal(t, worker.WorkerType(), WorkerTypeScraper)
}

func TestScraperWorkerExecute(t *testing.T) {
	var (
		mockScraper           = scraper.NewMockScraper()
		mockCoordinatorClient = coordinator_client.NewMockCoordinatorClient()
		scraperWorker         = NewScraperWorker(mockScraper, mockCoordinatorClient)
	)

	mockUrlTask, err := coordinator_client.NewTask("id", "test", ScraperWorkerParams{Url: "https://example.com"})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	mockScraper.SetHtmlContent("https://example.com", "<html><body><main>Hello, world!</main></body></html>")

	err = scraperWorker.Execute(context.Background(), mockUrlTask)
	assert.NoError(t, err)

	createdRagTask, err := mockCoordinatorClient.GetTask(context.Background(), time.Second*1, coordinator_client.CoordinatorClientTaskTopicRag)

	assert.NoError(t, err)

	parsedRagParams, err := coordinator_client.CastParams[RagWorkerParams](createdRagTask.Params)
	assert.NoError(t, err)
	assert.Equal(t, parsedRagParams.Markdown, "Hello, world!")
	assert.Equal(t, parsedRagParams.Url, "https://example.com")
}

func TestScraperWorkerExecuteNoMarkdown(t *testing.T) {
	var (
		mockScraper           = scraper.NewMockScraper()
		mockCoordinatorClient = coordinator_client.NewMockCoordinatorClient()
		scraperWorker         = NewScraperWorker(mockScraper, mockCoordinatorClient)
	)

	mockScraper.SetHtmlContent("https://example.com", "<html><body>Hello, world!</body></html>")

	mockUrlTask, err := coordinator_client.NewTask("id", "test", ScraperWorkerParams{Url: "https://example.com"})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = scraperWorker.Execute(context.Background(), mockUrlTask)
	assert.NoError(t, err)

	_, err = mockCoordinatorClient.GetTask(context.Background(), time.Second*1, coordinator_client.CoordinatorClientTaskTopicRag)
	if err != nil && err != coordinator_client.ErrNoTasksToComplete {
		t.Fatalf("Expected ErrNoTasksToComplete, got %v", err)
	}
}
