package worker

import (
	"testing"

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
