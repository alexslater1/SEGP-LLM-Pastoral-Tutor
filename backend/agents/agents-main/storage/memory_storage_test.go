package storage

import (
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("store and get single item", func(t *testing.T) {
		data := map[string]interface{}{
			"id":   "1",
			"name": "test",
		}

		// Test store
		_, err := storage.store(TableNameAgents, data)
		if err != nil {
			t.Errorf("Failed to store: %v", err)
		}

		// Test get
		retrieved, err := storage.get(TableNameAgents, "1")
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
		stored, err := storage.storeAll(TableNameAgents, data)
		if err != nil {
			t.Errorf("Failed to store multiple: %v", err)
		}
		if len(stored) != 2 {
			t.Errorf("Expected 2 items stored, got %d", len(stored))
		}

		// Test getAll with no filter
		retrieved, err := storage.getAll(TableNameAgents, map[string]string{})
		if err != nil {
			t.Errorf("Failed to get all: %v", err)
		}
		if len(retrieved) != 3 { // Including the previous test's item
			t.Errorf("Expected 3 items retrieved, got %d", len(retrieved))
		}
	})

	t.Run("getAll with filter", func(t *testing.T) {
		retrieved, err := storage.getAll(TableNameAgents, map[string]string{"name": "test2"})
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
		_, err := storage.get(TableNameAgents, "nonexistent")
		if err == nil {
			t.Error("Expected error when getting non-existent item")
		}
	})

	t.Run("store invalid data", func(t *testing.T) {
		invalidData := map[string]interface{}{
			"name": "test", // Missing ID
		}
		_, err := storage.store(TableNameAgents, invalidData)
		if err == nil {
			t.Error("Expected error when storing data without ID")
		}
	})
}
