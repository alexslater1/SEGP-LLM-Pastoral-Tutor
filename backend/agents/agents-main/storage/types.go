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
	UserID    string      `json:"user_id"`
}

func (ar AgentRequest) TableName() StorageTableName {
	return StorageTableNameAgentRequests
}

func NewAgentRequest(endpoint string, metadata interface{}, userID string) AgentRequest {
	return AgentRequest{
		Endpoint: endpoint,
		Metadata: metadata,
		UserID:   userID,
	}
}

type AgentEvent struct {
	ID        int         `json:"id,omitempty"`
	CreatedAt *time.Time  `json:"created_at,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Type      string      `json:"type"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

func (ae AgentEvent) TableName() StorageTableName {
	return StorageTableNameAgentEvents
}

func NewAgentEvent(requestId string, eventType string, metadata interface{}) AgentEvent {
	return AgentEvent{
		RequestID: requestId,
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

type Session struct {
	ID                 string     `json:"id,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	CreatedByRequestId string     `json:"created_by_request_id"`
	Name               *string    `json:"name,omitempty"`
	UserID             string     `json:"user_id"`
}

func (s Session) TableName() StorageTableName {
	return StorageTableNameSessions
}

func NewSession(createdByRequestId string, userID string) Session {
	return Session{
		CreatedByRequestId: createdByRequestId,
		UserID:             userID,
	}
}

type RequestSession struct {
	ID        string     `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	SessionID string     `json:"session_id"`
	RequestID string     `json:"request_id"`
}

func NewRequestSession(sessionID string, requestID string) RequestSession {
	return RequestSession{
		SessionID: sessionID,
		RequestID: requestID,
	}
}

func (rs RequestSession) TableName() StorageTableName {
	return StorageTableNameRequestSessions
}

type ChatCheck struct {
	ID        int        `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func NewChatCheck() ChatCheck {
	return ChatCheck{}
}

func (ls ChatCheck) TableName() StorageTableName {
	return StorageTableNameChatChecks
}

type CompletionResult struct {
	ID        int        `json:"id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	RequestID string     `json:"request_id"`
	Result    *string    `json:"result,omitempty"`
	Reason    *string    `json:"reason,omitempty"`
	Error     *string    `json:"error,omitempty"`
}

func (cr CompletionResult) TableName() StorageTableName {
	return StorageTableNameCompletionResults
}

func NewCompletionResult(requestID string, result *string, reason *string, err error) CompletionResult {
	if err != nil {
		errStr := err.Error()
		return CompletionResult{
			RequestID: requestID,
			Result:    result,
			Reason:    reason,
			Error:     &errStr,
		}
	}

	return CompletionResult{
		RequestID: requestID,
		Result:    result,
		Reason:    reason,
		Error:     nil,
	}
}
