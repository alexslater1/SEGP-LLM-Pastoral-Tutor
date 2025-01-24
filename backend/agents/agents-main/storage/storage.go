package storage

type StorageType interface {
	TableName() TableName
}

type TableName string

const (
	TableNameAgents TableName = "agents"
)

type Storage interface {
	store(table TableName, data interface{}) (interface{}, error)
	storeAll(table TableName, data []interface{}) ([]interface{}, error)

	get(table TableName, id string) (interface{}, error)
	getAll(table TableName, matchingFields map[string]string) ([]interface{}, error)
}
