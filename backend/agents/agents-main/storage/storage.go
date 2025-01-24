package storage

import (
	"fmt"
)

type StorageTableName string

const (
	StorageTableNameAgentRequests StorageTableName = "agent_requests"
)

type Storage interface {
	store(table StorageTableName, data interface{}) (interface{}, error)
	storeAll(table StorageTableName, data []interface{}) ([]interface{}, error)

	get(table StorageTableName, id string) (interface{}, error)
	getAll(table StorageTableName, matchingFields map[string]string) ([]interface{}, error)
}

func Get[T StorageType](storage Storage, id string) (T, error) {
	var t T
	data, err := storage.get(t.TableName(), id)
	if err != nil {
		return t, err
	}

	resultT, ok := data.(T)
	if !ok {
		return t, fmt.Errorf("failed to convert result to required type")
	}
	return resultT, nil
}

func GetAll[T StorageType](storage Storage, matchingFields map[string]string) ([]T, error) {
	var t T
	data, err := storage.getAll(t.TableName(), matchingFields)
	if err != nil {
		return nil, err
	}

	result := make([]T, len(data))
	for i, d := range data {
		resultT, ok := d.(T)
		if !ok {
			return nil, fmt.Errorf("failed to convert result to required type")
		}
		result[i] = resultT
	}

	return result, nil
}

func Store[T StorageType](storage Storage, data T) (T, error) {
	var t T
	result, err := storage.store(t.TableName(), data)
	if err != nil {
		return t, err
	}

	resultT, ok := result.(T)
	if !ok {
		return t, fmt.Errorf("failed to convert result to required type")
	}

	return resultT, nil
}

func StoreAll[T StorageType](storage Storage, data []T) ([]T, error) {
	var t T
	interfaceSlice := make([]interface{}, len(data))
	for i, v := range data {
		interfaceSlice[i] = v
	}
	result, err := storage.storeAll(t.TableName(), interfaceSlice)
	if err != nil {
		return nil, err
	}

	resultT := make([]T, len(result))
	for i, r := range result {
		rt, ok := r.(T)
		if !ok {
			return nil, fmt.Errorf("failed to convert result to required type")
		}
		resultT[i] = rt
	}

	return resultT, nil
}
