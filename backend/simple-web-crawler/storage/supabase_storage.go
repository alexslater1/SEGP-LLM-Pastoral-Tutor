package storage

import (
	"errors"

	supa "github.com/nedpals/supabase-go"
	postgrest_go "github.com/nedpals/supabase-go/postgrest/pkg"
)

type SupabaseStorage struct {
	client *supa.Client
}

func NewSupabaseStorage(supabaseUrl, supabaseServiceKey string) *SupabaseStorage {

	return &SupabaseStorage{
		client: supa.CreateClient(supabaseUrl, supabaseServiceKey),
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

	return result[0], nil
}

func (s *SupabaseStorage) storeAll(table StorageTableName, data []interface{}) ([]interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Insert(data).Execute(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
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

	return result[0], nil
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

	if filterQuery != nil {
		err := filterQuery.Execute(&results)
		return results, err
	}

	err := query.Execute(&results)
	return results, err
}

func (s *SupabaseStorage) update(table StorageTableName, id string, updateFields map[string]interface{}) (interface{}, error) {
	var results []interface{}
	err := s.client.DB.From(string(table)).Update(updateFields).Eq("id", id).Execute(&results)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no result returned from supabase")
	}

	return results[0], nil
}
