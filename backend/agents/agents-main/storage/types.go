package storage

type StorageType interface {
	TableName() StorageTableName
}

type AgentRequest struct {
	ID       *string `json:"id,omitempty"`
	Endpoint string  `json:"endpoint"`
}
