package storage

import (
	"fmt"
	"sync"
)

type MemoryStorage struct {
	data map[TableName]map[string]interface{}
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[TableName]map[string]interface{}),
	}
}

func (s *MemoryStorage) store(table TableName, data interface{}) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[table] == nil {
		s.data[table] = make(map[string]interface{})
	}

	// Assuming data has an ID field
	record, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("data must be a map[string]interface{}")
	}

	id, ok := record["id"].(string)
	if !ok {
		return nil, fmt.Errorf("data must have an 'id' field of type string")
	}

	s.data[table][id] = data
	return data, nil
}

func (s *MemoryStorage) storeAll(table TableName, data []interface{}) ([]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[table] == nil {
		s.data[table] = make(map[string]interface{})
	}

	result := make([]interface{}, len(data))
	for i, item := range data {
		record, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("data must be a map[string]interface{}")
		}

		id, ok := record["id"].(string)
		if !ok {
			return nil, fmt.Errorf("data must have an 'id' field of type string")
		}

		s.data[table][id] = item
		result[i] = item
	}

	return result, nil
}

func (s *MemoryStorage) get(table TableName, id string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data[table] == nil {
		return nil, fmt.Errorf("table %s not found", table)
	}

	if data, exists := s.data[table][id]; exists {
		return data, nil
	}

	return nil, fmt.Errorf("record with id %s not found in table %s", id, table)
}

func (s *MemoryStorage) getAll(table TableName, matchingFields map[string]string) ([]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data[table] == nil {
		return []interface{}{}, nil
	}

	var result []interface{}
	for _, record := range s.data[table] {
		matches := true
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		for field, value := range matchingFields {
			if recordValue, exists := recordMap[field]; !exists || recordValue != value {
				matches = false
				break
			}
		}

		if matches {
			result = append(result, record)
		}
	}

	return result, nil
}
