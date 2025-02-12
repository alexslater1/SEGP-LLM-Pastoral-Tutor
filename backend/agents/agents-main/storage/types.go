package storage

import (
	"time"
)

type StorageType interface {
	TableName() StorageTableName
}

type AgentRequest struct {
	ID        string      `json:"id,omitempty"`
	CreatedAt *time.Time  `json:"created_at,omitempty"`
	Endpoint  string      `json:"endpoint"`
	Metadata  interface{} `json:"metadata,omitempty"`
	ChatID    string      `json:"chat_id"`
}

func (ar AgentRequest) TableName() StorageTableName {
	return StorageTableNameAgentRequests
}

func NewAgentRequest(endpoint string, metadata interface{}, chatID string) AgentRequest {
	return AgentRequest{
		Endpoint: endpoint,
		Metadata: metadata,
		ChatID:   chatID,
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

type Chat struct {
	ID         string     `json:"id,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UserID     string     `json:"userId"`
	Title      *string    `json:"title,omitempty"`
	Visibility string     `json:"visibility"`
}

func (c Chat) TableName() StorageTableName {
	return StorageTableNameChats
}

func NewChat(userId string, title *string) Chat {
	return Chat{
		UserID:     userId,
		Visibility: "private",
		Title:      title,
	}
}
