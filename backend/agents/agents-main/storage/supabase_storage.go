package storage

import (
	supa "github.com/nedpals/supabase-go"
	postgrest_go "github.com/nedpals/supabase-go/postgrest/pkg"
	"github.com/segp/agents-main/utils"
)

type SupabaseStorage struct {
	client *supa.Client
}

func NewSupabaseStorage(supabaseUrl, supabaseServiceKey string) *SupabaseStorage {

	return &SupabaseStorage{
		client: supa.CreateClient(utils.Required(supabaseUrl, "supabaseUrl"), utils.Required(supabaseServiceKey, "supabaseServiceKey")),
	}
}

func (s *SupabaseStorage) store(table TableName, data interface{}) (interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Insert(data).Execute(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SupabaseStorage) storeAll(table TableName, data []interface{}) ([]interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Insert(data).Execute(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SupabaseStorage) get(table TableName, id string) (interface{}, error) {
	var result []interface{}
	err := s.client.DB.From(string(table)).Select("*").Eq("id", id).Execute(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SupabaseStorage) getAll(table TableName, matchingFields map[string]string) ([]interface{}, error) {
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

	return results, err
}
