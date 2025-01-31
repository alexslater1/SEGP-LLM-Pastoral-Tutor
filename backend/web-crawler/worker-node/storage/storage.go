package storage

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type StorageTableName string

const (
	StorageTableNameAgentRequests StorageTableName = "agent_requests"
	StorageTableNameAgentEvents   StorageTableName = "agent_events"
	StorageTableNameRagChunks     StorageTableName = "rag_chunks"
	StorageTableNameRagSources    StorageTableName = "rag_sources"
	StorageTableNameRagContacts   StorageTableName = "rag_contacts"
)

type Storage interface {
	store(table StorageTableName, data interface{}) (interface{}, error)
	storeAll(table StorageTableName, data []interface{}) ([]interface{}, error)

	get(table StorageTableName, id string) (interface{}, error)
	getAll(table StorageTableName, matchingFields map[string]string) ([]interface{}, error)
}

func Get[T StorageType](storage Storage, id string) (*T, error) {
	var t T

	data, err := storage.get(t.TableName(), id)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	ret := new(T)
	err = json.Unmarshal(jsonData, ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data into type %v: %v", reflect.TypeOf(t), err)
	}

	return ret, nil
}

func GetAll[T StorageType](storage Storage, matchingFields map[string]string) ([]T, error) {
	var t T
	data, err := storage.getAll(t.TableName(), matchingFields)
	if err != nil {
		return nil, err
	}

	ret := make([]T, len(data))
	for i, d := range data {
		jsonData, err := json.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
		}

		err = json.Unmarshal(jsonData, &ret[i])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data into type %v: %v", reflect.TypeOf(t), err)
		}
	}

	return ret, nil
}

func Store[T StorageType](storage Storage, data T) (*T, error) {
	var t T

	d, err := storage.store(t.TableName(), data)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	ret := new(T)
	err = json.Unmarshal(jsonData, ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data into type %v: %v", reflect.TypeOf(t), err)
	}

	return ret, nil
}

func StoreAll[T StorageType](storage Storage, data ...T) ([]T, error) {
	var t T

	table := t.TableName()

	// Convert []T to []interface{}
	converted := make([]interface{}, len(data))
	for i, v := range data {
		converted[i] = v
	}

	ds, err := storage.storeAll(table, converted)
	if err != nil {
		return nil, err
	}

	ret := make([]T, len(ds))
	for i, d := range ds {
		jsonData, err := json.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data to JSON: %v", err)
		}

		err = json.Unmarshal(jsonData, &ret[i])
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data into type %v: %v", reflect.TypeOf(t), err)
		}
	}

	return ret, nil
}
