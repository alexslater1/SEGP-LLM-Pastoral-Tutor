package storage

import (
	"testing"
)

type TestStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TestStructWithDifferentTag struct {
	ID   string `json:"id"`
	Name string `json:"different_name"`
}

type TestStructWithoutTags struct {
	ID   string
	Name string
}

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("store and get single item", func(t *testing.T) {
		data := map[string]interface{}{
			"id":   "1",
			"name": "test",
		}

		// Test store
		_, err := storage.store(StorageTableNameAgentRequests, data)
		if err != nil {
			t.Errorf("Failed to store: %v", err)
		}

		// Test get
		retrieved, err := storage.get(StorageTableNameAgentRequests, "1")
		if err != nil {
			t.Errorf("Failed to get: %v", err)
		}

		retrievedMap := retrieved.(map[string]interface{})
		if retrievedMap["id"] != "1" || retrievedMap["name"] != "test" {
			t.Errorf("Retrieved data doesn't match stored data")
		}
	})

	t.Run("store and get multiple items", func(t *testing.T) {
		data := []interface{}{
			map[string]interface{}{
				"id":   "2",
				"name": "test2",
			},
			map[string]interface{}{
				"id":   "3",
				"name": "test3",
			},
		}

		// Test storeAll
		stored, err := storage.storeAll(StorageTableNameAgentRequests, data)
		if err != nil {
			t.Errorf("Failed to store multiple: %v", err)
		}
		if len(stored) != 2 {
			t.Errorf("Expected 2 items stored, got %d", len(stored))
		}

		// Test getAll with no filter
		retrieved, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{})
		if err != nil {
			t.Errorf("Failed to get all: %v", err)
		}
		if len(retrieved) != 3 { // Including the previous test's item
			t.Errorf("Expected 3 items retrieved, got %d", len(retrieved))
		}
	})

	t.Run("getAll with filter", func(t *testing.T) {
		retrieved, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{"name": "test2"})
		if err != nil {
			t.Errorf("Failed to get filtered: %v", err)
		}
		if len(retrieved) != 1 {
			t.Errorf("Expected 1 filtered item, got %d", len(retrieved))
		}

		retrievedMap := retrieved[0].(map[string]interface{})
		if retrievedMap["name"] != "test2" {
			t.Errorf("Retrieved filtered data doesn't match expected")
		}
	})

	t.Run("get non-existent item", func(t *testing.T) {
		_, err := storage.get(StorageTableNameAgentRequests, "nonexistent")
		if err == nil {
			t.Error("Expected error when getting non-existent item")
		}
	})

	t.Run("store invalid data", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"name": "test", // Missing ID
		}
		_, err := storage.store(StorageTableNameAgentRequests, invalidData)
		if err == nil {
			t.Error("Expected error when storing data without ID")
		}
	})
}

func TestMemoryStorageWithStructs(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("store and get struct", func(t *testing.T) {
		data := TestStruct{
			ID:   "struct1",
			Name: "test struct",
		}

		// Test store
		_, err := storage.store(StorageTableNameAgentRequests, data)
		if err != nil {
			t.Errorf("Failed to store struct: %v", err)
		}

		// Test get
		retrieved, err := storage.get(StorageTableNameAgentRequests, "struct1")
		if err != nil {
			t.Errorf("Failed to get struct: %v", err)
		}

		// Convert retrieved data to map for comparison
		retrievedMap, err := structToMap(retrieved)
		if err != nil {
			t.Errorf("Failed to convert retrieved data to map: %v", err)
		}

		if retrievedMap["id"] != "struct1" || retrievedMap["name"] != "test struct" {
			t.Errorf("Retrieved struct data doesn't match stored data")
		}
	})

	t.Run("store and get multiple structs", func(t *testing.T) {
		data := []interface{}{
			TestStruct{
				ID:   "struct2",
				Name: "test struct 2",
			},
			TestStruct{
				ID:   "struct3",
				Name: "test struct 3",
			},
		}

		// Test storeAll
		stored, err := storage.storeAll(StorageTableNameAgentRequests, data)
		if err != nil {
			t.Errorf("Failed to store multiple structs: %v", err)
		}
		if len(stored) != 2 {
			t.Errorf("Expected 2 structs stored, got %d", len(stored))
		}

		// Test getAll with name filter
		retrieved, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{"name": "test struct 2"})
		if err != nil {
			t.Errorf("Failed to get filtered structs: %v", err)
		}
		if len(retrieved) != 1 {
			t.Errorf("Expected 1 filtered struct, got %d", len(retrieved))
		}
	})

	t.Run("store struct with different json tags", func(t *testing.T) {
		data := TestStructWithDifferentTag{
			ID:   "struct4",
			Name: "test different tags",
		}

		// Test store
		_, err := storage.store(StorageTableNameAgentRequests, data)
		if err != nil {
			t.Errorf("Failed to store struct with different tags: %v", err)
		}

		// Test get and verify the field mapping
		retrieved, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{"different_name": "test different tags"})
		if err != nil {
			t.Errorf("Failed to get struct with different tags: %v", err)
		}
		if len(retrieved) != 1 {
			t.Errorf("Expected 1 struct with matching different_name, got %d", len(retrieved))
		}
	})

	t.Run("store struct without json tags", func(t *testing.T) {
		data := TestStructWithoutTags{
			ID:   "struct5",
			Name: "test no tags",
		}

		// Test store - we expect this to fail
		_, err := storage.store(StorageTableNameAgentRequests, data)
		if err == nil {
			t.Error("Expected error when storing struct without JSON tags")
		}

		// Since store should fail, we don't need to test retrieval
	})

	t.Run("store invalid struct without ID", func(t *testing.T) {
		type InvalidStruct struct {
			Name string `json:"name"`
		}

		data := InvalidStruct{
			Name: "test invalid",
		}

		_, err := storage.store(StorageTableNameAgentRequests, data)
		if err == nil {
			t.Error("Expected error when storing struct without ID field")
		}
	})
}
