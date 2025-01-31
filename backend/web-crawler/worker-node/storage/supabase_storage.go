package storage

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethanhosier/worker-node/utils"
	supa "github.com/nedpals/supabase-go"
	postgrest_go "github.com/nedpals/supabase-go/postgrest/pkg"
)

type SupabaseStorage struct {
	client *supa.Client
}

func NewSupabaseStorage(supabaseUrl, supabaseServiceKey string) *SupabaseStorage {

	return &SupabaseStorage{
		client: supa.CreateClient(utils.Required(supabaseUrl, "supabaseUrl"), utils.Required(supabaseServiceKey, "supabaseServiceKey")),
	}
}

func (s *SupabaseStorage) store(table StorageTableName, data interface{}) (interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Insert(data).Execute(&result)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no result returned from supabase")
	}

	return processEmbeddingField(result[0])
}

func (s *SupabaseStorage) storeAll(table StorageTableName, data []interface{}) ([]interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Insert(data).Execute(&result)

	if err != nil {
		return nil, err
	}

	return processAllEmbeddingFields(result)
}

func (s *SupabaseStorage) get(table StorageTableName, id string) (interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Select("*").Eq("id", id).Execute(&result)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no result returned from supabase")
	}

	return processEmbeddingField(result[0])
}

func (s *SupabaseStorage) getAll(table StorageTableName, matchingFields map[string]string) ([]interface{}, error) {
	var results []interface{}

	query := s.client.DB.From(string(table)).Select("*")
	var filterQuery *postgrest_go.FilterRequestBuilder
	for k, v := range matchingFields {
		if filterQuery == nil {
			filterQuery = query.Filter(k, "eq", v)
		} else {
			filterQuery = filterQuery.Filter(k, "eq", v)
		}
	}
	err := filterQuery.Execute(&results)
	if err != nil {
		return nil, err
	}

	return processAllEmbeddingFields(results)
}

func processAllEmbeddingFields(data []interface{}) ([]interface{}, error) {
	for i, d := range data {
		processed, err := processEmbeddingField(d)
		if err != nil {
			return nil, err
		}
		data[i] = processed
	}
	return data, nil
}

func processEmbeddingField(data interface{}) (interface{}, error) {
	// Convert interface{} to map[string]interface{}
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return data, nil // Return original if not a map
	}

	// Check if embedding field exists
	embeddingRaw, exists := dataMap["embedding"]
	if !exists {
		return data, nil // Return original if no embedding field
	}

	// Convert embedding string to []float32
	embeddingStr, ok := embeddingRaw.(string)
	if !ok {
		return nil, fmt.Errorf("embedding field is not a string")
	}

	var embedding []float32
	err := parseStringToFloat32Slice(embeddingStr, &embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedding: %w", err)
	}

	// Update the map with the converted embedding
	dataMap["embedding"] = embedding
	return dataMap, nil
}

func parseStringToFloat32Slice(s string, result *[]float32) error {
	// First parse into []float64 since JSON unmarshaling doesn't directly support float32
	var temp []float64
	if err := json.Unmarshal([]byte(s), &temp); err != nil {
		return fmt.Errorf("failed to unmarshal string to float slice: %w", err)
	}

	// Convert float64 to float32
	*result = make([]float32, len(temp))
	for i, v := range temp {
		(*result)[i] = float32(v)
	}

	return nil
}
