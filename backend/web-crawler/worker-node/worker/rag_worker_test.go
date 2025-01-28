package worker

import (
	"context"
	"testing"
	"time"

	"github.com/ethanhosier/worker-node/coordinator_client"
	"github.com/ethanhosier/worker-node/ragger"
	"github.com/ethanhosier/worker-node/storage"
	"github.com/stretchr/testify/assert"
)

func TestRagWorkerWorkerType(t *testing.T) {
	mock := NewRagWorker(nil, nil, nil)
	if mock.WorkerType() != WorkerTypeRag {
		t.Errorf("WorkerType should be WorkerTypeRag")
	}
}

func TestRagWorkerId(t *testing.T) {
	ragWorker := NewRagWorker(nil, nil, nil)
	if ragWorker.Id() == "" {
		t.Errorf("Id should not be empty")
	}
}

func TestStoreWebsite(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragWorker     = NewRagWorker(nil, nil, memoryStorage)
		url           = "https://example.com"
	)

	storedWebsite, err := ragWorker.storeWebsite(url)
	if err != nil {
		t.Errorf("Error storing website: %v", err)
	}

	assert.Equal(t, storedWebsite.URL, url)
}

func TestRagWorkerProcessAndStoreChunks(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragClient     = ragger.NewMockRagClient()
		ragWorker     = NewRagWorker(ragClient, nil, memoryStorage)
		chunks        = []string{"Hello, world!1", "Hello, world!2", "Hello, world!3"}
		embeddings    = [][]float32{{1.0, 2.0, 3.0}, {4.0, 5.0, 6.0}, {7.0, 8.0, 9.0}}
	)

	ragClient.SetChunksFor("Hello, world!", chunks)
	for i, embedding := range embeddings {
		ragClient.SetEmbeddingsFor(chunks[i], embedding)
	}

	ragWorker.processAndStoreChunks("Hello, world!", 1)

	rags, err := storage.GetAll[storage.Rag](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting chunks: %v", err)
	}

	for _, rag := range rags {
		assert.Equal(t, rag.Text, chunks[rag.PosInDoc])
		assert.Equal(t, rag.Embedding, embeddings[rag.PosInDoc])
		assert.Equal(t, rag.WebsiteID, 1)
	}
}

func TestRagWorkerProcessAndStoreContacts(t *testing.T) {
	var (
		memoryStorage = storage.NewMemoryStorage()
		ragClient     = ragger.NewMockRagClient()
		ragWorker     = NewRagWorker(ragClient, nil, memoryStorage)
	)

	ragClient.SetContactsFor("Hello, world!", []ragger.Contact{
		{Context: "Hello, world!", Value: "John Doe", Type: "person"},
	})

	ragWorker.processAndStoreContacts("Hello, world!", 1)

	contacts, err := storage.GetAll[storage.Contacts](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting contacts: %v", err)
	}

	assert.Equal(t, len(contacts), 1)
	assert.Equal(t, contacts[0].Context, "Hello, world!")
	assert.Equal(t, contacts[0].Contact, "John Doe")
	assert.Equal(t, contacts[0].ContactType, "person")
	assert.Equal(t, contacts[0].DocID, 1)
}

func TestRagWorkerExecute(t *testing.T) {
	// given
	var (
		memoryStorage     = storage.NewMemoryStorage()
		ragClient         = ragger.NewMockRagClient()
		coordinatorClient = coordinator_client.NewMockCoordinatorClient()
		ragWorker         = NewRagWorker(ragClient, coordinatorClient, memoryStorage)

		websiteUrl = "https://example.com"
		markdown   = "Hello, world!"

		chunks     = []string{"Hello, world!1"}
		embeddings = [][]float32{{1.0, 2.0, 3.0}}

		contacts = []ragger.Contact{{Context: markdown, Value: "John Doe", Type: "person"}}
	)

	task, err := coordinator_client.NewTask("1", "test", RagWorkerParams{
		Markdown: markdown,
		Url:      websiteUrl,
	})
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}

	ragClient.SetChunksFor(markdown, chunks)
	ragClient.SetContactsFor(markdown, contacts)
	for i, embedding := range embeddings {
		ragClient.SetEmbeddingsFor(chunks[i], embedding)
	}

	// when
	err = ragWorker.Execute(context.TODO(), task)
	if err != nil {
		t.Errorf("Error executing task: %v", err)
	}

	// then
	websites, err := storage.GetAll[storage.Website](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting websites: %v", err)
	}

	assert.Equal(t, len(websites), 1)
	assert.Equal(t, websites[0].URL, websiteUrl)

	rags, err := storage.GetAll[storage.Rag](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting chunks: %v", err)
	}

	assert.Equal(t, len(rags), 1)
	assert.Equal(t, rags[0].Text, chunks[0])
	assert.Equal(t, rags[0].Embedding, embeddings[0])
	assert.Equal(t, rags[0].WebsiteID, websites[0].ID)

	storedContacts, err := storage.GetAll[storage.Contacts](memoryStorage, nil)
	if err != nil {
		t.Errorf("Error getting contacts: %v", err)
	}

	assert.Equal(t, len(storedContacts), 1)
	assert.Equal(t, storedContacts[0].Context, markdown)
	assert.Equal(t, storedContacts[0].Contact, "John Doe")
	assert.Equal(t, storedContacts[0].ContactType, "person")
	assert.Equal(t, storedContacts[0].DocID, websites[0].ID)
}

func TestRagWorkerCleanup(t *testing.T) {
	// given
	var (
		coordinatorClient = coordinator_client.NewMockCoordinatorClient()
		ragWorker         = NewRagWorker(nil, coordinatorClient, nil)
	)

	task, err := coordinator_client.NewTask("1", "test", RagWorkerParams{
		Markdown: "Hello, world!",
		Url:      "https://example.com",
	})
	if err != nil {
		t.Errorf("Error creating task: %v", err)
	}

	coordinatorClient.CreateTask(context.TODO(), coordinator_client.CoordinatorClientTaskTopicRag, task)
	if _, err := coordinatorClient.GetTaskAndSetProcessing(context.TODO(), 1*time.Second, coordinator_client.CoordinatorClientTaskTopicRag); err != nil {
		t.Errorf("Error getting task: %v", err)
	}

	// when
	err = ragWorker.Cleanup(context.TODO(), task)

	// then
	if err != nil {
		t.Errorf("Error cleaning up: %v", err)
	}

	if err := coordinatorClient.SetProcessed(context.TODO(), coordinator_client.CoordinatorClientTaskTopicRag, task); err != coordinator_client.ErrNoTasksCompleted {
		t.Errorf("There should be an error setting task to processed: %v", err)
	}
}
