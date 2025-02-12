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
	ChatID   string      `json:"chat_id,omitempty"`
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
	UserID     string     `json:"userId,omitempty"`
	Title      string     `json:"title,omitempty"`
	Visibility string     `json:"visibility,omitempty"`
}

func (c Chat) TableName() StorageTableName {
	return StorageTableNameChats
}

func NewChat(id string, userId string, title string, visibility string) Chat {
	return Chat{
		ID:         id,
		UserID:     userId,
		Title:      title,
		Visibility: visibility,
	}
}
