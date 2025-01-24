package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Implement TableName for AgentRequest
func (ar AgentRequest) TableName() StorageTableName {
	return StorageTableNameAgentRequests
}

func TestStore(t *testing.T) {
	storage := NewMemoryStorage()

	request := AgentRequest{
		ID:       stringPtr("1"),
		Endpoint: "/test",
	}

	stored, err := Store(storage, request)
	assert.NoError(t, err)
	assert.Equal(t, request, stored)

	// Test storing without ID
	invalidRequest := AgentRequest{
		Endpoint: "/test",
	}
	_, err = Store(storage, invalidRequest)
	assert.Error(t, err)
}

func TestStoreAll(t *testing.T) {
	storage := NewMemoryStorage()

	requests := []AgentRequest{
		{
			ID:       stringPtr("1"),
			Endpoint: "/test1",
		},
		{
			ID:       stringPtr("2"),
			Endpoint: "/test2",
		},
	}

	stored, err := StoreAll(storage, requests)
	assert.NoError(t, err)
	assert.Equal(t, requests, stored)

	// Test storing invalid requests
	invalidRequests := []AgentRequest{
		{
			ID:       stringPtr("3"),
			Endpoint: "/test3",
		},
		{
			Endpoint: "/test4", // Missing ID
		},
	}
	_, err = StoreAll(storage, invalidRequests)
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	storage := NewMemoryStorage()

	request := AgentRequest{
		ID:       stringPtr("1"),
		Endpoint: "/test",
	}

	// Store first
	_, err := Store(storage, request)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := Get[AgentRequest](storage, "1")
	assert.NoError(t, err)
	assert.Equal(t, request, retrieved)

	// Test Get with non-existent ID
	_, err = Get[AgentRequest](storage, "999")
	assert.Error(t, err)
}

func TestGetAll(t *testing.T) {
	storage := NewMemoryStorage()

	requests := []AgentRequest{
		{
			ID:       stringPtr("1"),
			Endpoint: "/test1",
		},
		{
			ID:       stringPtr("2"),
			Endpoint: "/test1", // Same endpoint
		},
		{
			ID:       stringPtr("3"),
			Endpoint: "/test2",
		},
	}

	// Store all requests
	_, err := StoreAll(storage, requests)
	assert.NoError(t, err)

	// Test GetAll with matching endpoint
	matchingFields := map[string]string{"endpoint": "/test1"}
	retrieved, err := GetAll[AgentRequest](storage, matchingFields)
	assert.NoError(t, err)
	assert.Len(t, retrieved, 2) // Should find 2 requests with endpoint "/test1"

	// Test GetAll with no matches
	noMatches := map[string]string{"endpoint": "/nonexistent"}
	retrieved, err = GetAll[AgentRequest](storage, noMatches)
	assert.NoError(t, err)
	assert.Empty(t, retrieved)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
