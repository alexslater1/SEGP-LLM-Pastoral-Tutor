package storage

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"

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

		// Preserve original integer type for ID field
		if originalID, ok := getOriginalIDType(data); ok && isNumberType(originalID) {
			if idFloat, ok := dataMap["id"].(float64); ok {
				dataMap["id"] = int(idFloat)
			}
		}
	}

	// If ID is empty, generate a random UUID
	if dataMap["id"] == nil {
		// Check if the original data had a numeric ID type
		if originalID, ok := getOriginalIDType(data); ok && isNumberType(originalID) {
			// Generate random number between 1 and 1000000 for numeric IDs
			dataMap["id"] = rand.Intn(1000000) + 1
		} else {
			// Default to UUID string if not numeric
			dataMap["id"] = uuid.New().String()
		}
	}

	// Get the ID as string or number
	id, ok := dataMap["id"].(string)
	if !ok {
		// Try as number
		if numID, ok := dataMap["id"].(int); ok {
			id = strconv.Itoa(numID)
		} else {
			return nil, fmt.Errorf("ID must be a string or number")
		}
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

func (s *MemoryStorage) getAll(table StorageTableName, query *QueryBuilder) ([]interface{}, error) {
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
		if query != nil {
			for field, value := range query.matchingFields {
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
		}

		if matches {
			result = append(result, item)
		}
	}

	// Sort the results based on query.orderBy
	if query != nil && query.orderBy != nil {
		sort.Slice(result, func(i, j int) bool {
			item1, _ := result[i].(map[string]interface{})
			item2, _ := result[j].(map[string]interface{})
			if query.orderBy.order == OrderByAsc {
				return item1[query.orderBy.column].(string) < item2[query.orderBy.column].(string)
			}
			return item1[query.orderBy.column].(string) > item2[query.orderBy.column].(string)
		})
	}

	r := result
	for field, value := range query.greaterThanFields {
		r = filterGt(r, field, value)
	}

	if query != nil && query.limit != nil {
		return r[:*query.limit], nil
	}

	return r, nil
}

// Helper function to check original ID type
func getOriginalIDType(data interface{}) (interface{}, bool) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		field := val.FieldByName("ID")
		if field.IsValid() {
			return field.Interface(), true
		}
	}
	if m, ok := data.(map[string]interface{}); ok {
		if id, exists := m["id"]; exists {
			return id, true
		}
	}
	return nil, false
}

// Helper function to check if type is numeric
func isNumberType(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	}
	return false
}

func (s *MemoryStorage) update(table StorageTableName, id string, updateFields map[string]interface{}) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.data[table][id]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}

	currentMap, ok := current.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("item is not a map")
	}

	for k, v := range updateFields {
		currentMap[k] = v
	}

	s.data[table][id] = currentMap
	return currentMap, nil
}

func (s *MemoryStorage) delete(table StorageTableName, id string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.data[table][id]
	if !ok {
		return nil, fmt.Errorf("item not found")
	}

	delete(s.data[table], id)

	return current, nil
}

func (s *MemoryStorage) deleteAll(table StorageTableName, query *QueryBuilder) ([]interface{}, error) {
	toDelete, err := s.getAll(table, query)
	if err != nil {
		return nil, err
	}

	for _, item := range toDelete {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("error casting type to map string interface")
		}

		s.delete(table, itemMap["id"].(string))
	}

	return toDelete, err
}

func filterGt(result []interface{}, fieldName string, value interface{}) []interface{} {
	filteredResult := []interface{}{}

	for _, item := range result {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemValue, ok := itemMap[fieldName]
		if !ok {
			continue
		}

		switch itemValue.(type) {
		case string:
			if itemValue.(string) > value.(string) {
				filteredResult = append(filteredResult, item)
			}
		case int:
			if itemValue.(int) > value.(int) {
				filteredResult = append(filteredResult, item)
			}
		case float64:
			if itemValue.(float64) > value.(float64) {
				filteredResult = append(filteredResult, item)
			}
		case time.Time:
			if itemValue.(time.Time).After(value.(time.Time)) {
				filteredResult = append(filteredResult, item)
			}
		}
	}

	return filteredResult
}
