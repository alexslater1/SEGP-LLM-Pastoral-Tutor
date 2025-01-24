package storage

type StorageType interface {
	TableName() StorageTableName
}

type AgentRequest struct {
	ID       string `json:"id,omitempty"`
	Endpoint string `json:"endpoint"`
}

func (ar AgentRequest) TableName() StorageTableName {
	return StorageTableNameAgentRequests
}

type AgentCompletionRequestEventType string

const (
	AgentCompletionRequestEventTypeThink   AgentCompletionRequestEventType = "Think"
	AgentCompletionRequestEventTypeAct     AgentCompletionRequestEventType = "Act"
	AgentCompletionRequestEventTypeObserve AgentCompletionRequestEventType = "Observe"
)

type AgentCompletionRequestEvent struct {
	ID             string      `json:"id,omitempty"`
	AgentRequestID string      `json:"agent_request_id"`
	Type           string      `json:"type"`
	Metadata       interface{} `json:"metadata"`
}

func (ac AgentCompletionRequestEvent) TableName() StorageTableName {
	return StorageTableNameAgentCompletionRequestEvents
}
