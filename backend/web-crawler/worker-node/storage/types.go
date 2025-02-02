package storage

import (
	"time"
)

type StorageType interface {
	TableName() StorageTableName
}

type AgentRequest struct {
	ID       string      `json:"id,omitempty"`
	Endpoint string      `json:"endpoint"`
	Metadata interface{} `json:"metadata,omitempty"`
}

func (ar AgentRequest) TableName() StorageTableName {
	return StorageTableNameAgentRequests
}

func NewAgentRequest(endpoint string, metadata interface{}) AgentRequest {
	return AgentRequest{
		Endpoint: endpoint,
		Metadata: metadata,
	}
}

type AgentEvent struct {
	ID        int         `json:"id,omitempty"`
	CreatedAt *time.Time  `json:"created_at,omitempty"`
	RequestId string      `json:"request_id"`
	Type      string      `json:"type"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

func (ae AgentEvent) TableName() StorageTableName {
	return StorageTableNameAgentEvents
}

func NewAgentEvent(requestId string, eventType string, metadata interface{}) AgentEvent {
	return AgentEvent{
		RequestId: requestId,
		Type:      eventType,
		Metadata:  metadata,
	}
}

type RagChunk struct {
	ID          int       `json:"id,omitempty"`
	RagSourceId int       `json:"rag_source_id"`
	Text        string    `json:"text"`
	PosInSource int       `json:"pos_in_source"`
	Embedding   []float32 `json:"embedding"`
}

func (r RagChunk) TableName() StorageTableName {
	return StorageTableNameRagChunks
}

type RagSource struct {
	ID   int    `json:"id,omitempty"`
	URL  string `json:"url"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r RagSource) TableName() StorageTableName {
	return StorageTableNameRagSources
}

type RagContact struct {
	ID          int       `json:"id,omitempty"`
	RagSourceId int       `json:"rag_source_id"`
	Context     string    `json:"context"`
	Contact     string    `json:"contact"`
	PosInSource int       `json:"pos_in_source"`
	ContactType string    `json:"contact_type"`
	Embedding   []float32 `json:"embedding"`
}

func (c RagContact) TableName() StorageTableName {
	return StorageTableNameRagContacts
}
