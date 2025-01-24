package storage

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	data map[StorageTableName]map[string]interface{} // Table -> ID -> Data
	mu   sync.RWMutex                                // For concurrent access
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[StorageTableName]map[string]interface{}),
	}
}

func (s *MemoryStorage) store(table StorageTableName, data interface{}) (interface{}, error) {
	// Check for nil data
	if data == nil {
		return nil, fmt.Errorf("data cannot be nil")
	}

	var dataMap map[string]interface{}

	// Try to convert data to map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		dataMap = v
	default:
		// Convert struct to map using reflection
		bytes, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal struct to json: %w", err)
		}

		if err := json.Unmarshal(bytes, &dataMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json to map: %w", err)
		}
	}

	// If ID is empty, generate a random UUID
	if dataMap["id"] == nil {
		dataMap["id"] = uuid.New().String()
	}

	// Get the ID as string
	id, ok := dataMap["id"].(string)
	if !ok {
		return nil, fmt.Errorf("ID must be a string")
	}

	// Initialize table if it doesn't exist
	if s.data[table] == nil {
		s.data[table] = make(map[string]interface{})
	}

	// Store the data
	s.data[table][id] = dataMap

	return dataMap, nil
}

func (s *MemoryStorage) storeAll(table StorageTableName, data []interface{}) ([]interface{}, error) {
	// Initialize result slice
	result := make([]interface{}, 0, len(data))

	// Process each item
	for _, item := range data {
		storedData, err := s.store(table, item)
		if err != nil {
			return nil, fmt.Errorf("failed to store item: %w", err)
		}
		result = append(result, storedData)
	}

	return result, nil
}

func (s *MemoryStorage) get(table StorageTableName, id string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[table][id]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}
	return item, nil
}

func (s *MemoryStorage) getAll(table StorageTableName, matchingFields map[string]string) ([]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []interface{}
	for _, item := range s.data[table] {
		// Convert item to map[string]interface{}
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if all matching fields match
		matches := true
		for field, value := range matchingFields {
			itemValue, exists := itemMap[field]
			if !exists {
				matches = false
				break
			}
			// Convert itemValue to string for comparison
			itemValueStr, ok := itemValue.(string)
			if !ok {
				matches = false
				break
			}
			if itemValueStr != value {
				matches = false
				break
			}
		}

		if matches {
			result = append(result, item)
		}
	}
	return result, nil
}
