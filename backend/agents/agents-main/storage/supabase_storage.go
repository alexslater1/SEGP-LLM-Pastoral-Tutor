package storage

import (
	"errors"
	"fmt"
	"time"

	supa "github.com/nedpals/supabase-go"
	postgrest_go "github.com/nedpals/supabase-go/postgrest/pkg"
	"github.com/segp/agents-main/utils"
)

type SupabaseStorage struct {
	client *supa.Client
}

func NewSupabaseStorage(supabaseUrl, supabaseServiceKey string) *SupabaseStorage {
	utils.Required(supabaseUrl, "supabaseUrl")
	utils.Required(supabaseServiceKey, "supabaseServiceKey")

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

func (s *SupabaseStorage) getAll(table StorageTableName, query *QueryBuilder) ([]interface{}, error) {
	return s.handleQuery(queryTypeSelect, table, query)
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

func (s *SupabaseStorage) delete(table StorageTableName, id string) (interface{}, error) {
	var results []interface{}
	err := s.client.DB.From(string(table)).Delete().Eq("id", id).Execute(&results)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no result returned from supabase")
	}

	return results[0], nil
}

func (s *SupabaseStorage) deleteAll(table StorageTableName, query *QueryBuilder) ([]interface{}, error) {
	return s.handleQuery(queryTypeDelete, table, query)
}

type queryType string

const (
	queryTypeSelect queryType = "select"
	queryTypeDelete queryType = "delete"
)

func (s *SupabaseStorage) handleQuery(queryType queryType, table StorageTableName, query *QueryBuilder) ([]interface{}, error) {
	queryBuilder := s.client.DB.From(string(table))

	switch queryType {
	case queryTypeSelect:
		return s.handleSelectQuery(queryBuilder, query)
	case queryTypeDelete:
		return s.handleDeleteQuery(queryBuilder, query)
	}

	panic("invalid query type")
}

func (s *SupabaseStorage) handleSelectQuery(requestBuilder *postgrest_go.RequestBuilder, query *QueryBuilder) ([]interface{}, error) {
	var results []interface{}
	selectRequest := requestBuilder.Select("*")

	if query == nil {
		err := selectRequest.Execute(&results)
		return results, err
	}

	for k, v := range query.matchingFields {
		selectRequest.Filter(k, "eq", v)
	}

	for k, v := range query.greaterThanFields {
		selectRequest.Filter(k, "gt", toString(v))
	}

	if query.orderBy != nil {
		selectRequest.OrderBy(query.orderBy.column, string(query.orderBy.order))
	}

	if query.limit != nil {
		selectRequest.Limit(*query.limit)
	}

	err := selectRequest.Execute(&results)
	return results, err
}

func (s *SupabaseStorage) handleDeleteQuery(requestBuilder *postgrest_go.RequestBuilder, query *QueryBuilder) ([]interface{}, error) {
	var results []interface{}

	if query == nil {
		err := requestBuilder.Delete().Filter("id", "neq", "a7517604-2f29-47af-855d-81a3808ddf67").Execute(&results)
		return results, err
	}

	deleteRequest := requestBuilder.Delete()

	for k, v := range query.matchingFields {
		deleteRequest.Filter(k, "eq", v)
	}

	for k, v := range query.greaterThanFields {
		deleteRequest.Filter(k, "gt", toString(v))
	}

	err := deleteRequest.Execute(&results)

	if err != nil {
		return nil, err
	}

	return results, nil
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return fmt.Sprintf("%+v", v.Format(time.RFC3339))
	}
	return fmt.Sprintf("%+v", value)
}
