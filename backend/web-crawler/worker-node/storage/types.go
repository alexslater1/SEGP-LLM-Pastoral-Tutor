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

type Rag struct {
	ID        int       `json:"id,omitempty"`
	Text      string    `json:"text"`
	DocID     int       `json:"doc_id"`
	PosInDoc  int       `json:"pos_in_doc"`
	Embedding []float32 `json:"embedding"`
	Source    string    `json:"source,omitempty"`
	WebsiteID int       `json:"website_id,omitempty"`
}

func (r Rag) TableName() StorageTableName {
	return StorageTableNameRag
}

type Website struct {
	ID  int    `json:"id,omitempty"`
	URL string `json:"url"`
}

func (w Website) TableName() StorageTableName {
	return StorageTableNameWebsite
}

type Contact struct {
	ID            int       `json:"id,omitempty"`
	Context       string    `json:"context"`
	DocID         int       `json:"doc_id"`
	Contact       string    `json:"contact"`
	PosInContacts int       `json:"pos_in_contacts"`
	ContactType   string    `json:"contact_type"`
	Embedding     []float32 `json:"embedding"`
	Source        string    `json:"source,omitempty"`
	WebsiteID     int       `json:"website_id,omitempty"`
}

func (c Contact) TableName() StorageTableName {
	return StorageTableNameContacts
}
