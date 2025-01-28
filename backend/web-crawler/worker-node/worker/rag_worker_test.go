package worker

import (
	"testing"

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
