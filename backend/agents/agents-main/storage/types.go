package storage

type StorageType interface {
	TableName() TableName
}
