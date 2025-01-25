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
