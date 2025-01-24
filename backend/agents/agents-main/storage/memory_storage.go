package storage

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type MemoryStorage struct {
	data map[StorageTableName]map[string]interface{}
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[StorageTableName]map[string]interface{}),
	}
}

func (s *MemoryStorage) store(table StorageTableName, data interface{}) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[table] == nil {
		s.data[table] = make(map[string]interface{})
	}

	// Convert data to map if needed
	var record map[string]interface{}
	switch v := data.(type) {
	case map[string]interface{}:
		record = v
	default:
		// Try to convert the struct to a map using reflection
		recordMap, err := structToMap(data)
		fmt.Printf("recordMap: %+v\n", recordMap)
		if err != nil {
			return nil, fmt.Errorf("failed to convert data to map: %w", err)
		}
		record = recordMap
	}

	id, ok := record["id"].(string)
	if !ok {
		return nil, fmt.Errorf("data must have an 'id' field of type string")
	}

	s.data[table][id] = data
	return data, nil
}

func (s *MemoryStorage) storeAll(table StorageTableName, data []interface{}) ([]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data[table] == nil {
		s.data[table] = make(map[string]interface{})
	}

	result := make([]interface{}, len(data))
	for i, item := range data {
		// Convert item to map if needed
		var record map[string]interface{}
		switch v := item.(type) {
		case map[string]interface{}:
			record = v
		default:
			// Try to convert the struct to a map using reflection
			recordMap, err := structToMap(item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert data to map: %w", err)
			}
			record = recordMap
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

func (s *MemoryStorage) get(table StorageTableName, id string) (interface{}, error) {
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

func (s *MemoryStorage) getAll(table StorageTableName, matchingFields map[string]string) ([]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data[table] == nil {
		return []interface{}{}, nil
	}

	var result []interface{}
	for _, record := range s.data[table] {
		matches := true

		// Convert the record to a map if it isn't already
		var recordMap map[string]interface{}
		switch v := record.(type) {
		case map[string]interface{}:
			recordMap = v
		default:
			converted, err := structToMap(record)
			if err != nil {
				continue
			}
			recordMap = converted
		}

		for field, value := range matchingFields {
			if recordValue, exists := recordMap[field]; !exists || fmt.Sprint(recordValue) != value {
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

// Helper function to convert struct to map using reflection
func structToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	// If it's a pointer, get the underlying element
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or a pointer to a struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// Dereference pointer fields
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
		}

		// Use the json tag if available, otherwise use the field name
		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}
		// Remove the omitempty option if present
		if idx := strings.Index(key, ","); idx != -1 {
			key = key[:idx]
		}
		if key == "-" {
			continue
		}
		result[key] = fieldValue.Interface()
	}

	return result, nil
}
