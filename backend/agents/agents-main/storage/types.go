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
	AgentCompletionRequestEventTypeError   AgentCompletionRequestEventType = "Error"
	AgentCompletionRequestEventTypeSuccess AgentCompletionRequestEventType = "Success"
)

type AgentCompletionRequestEvent struct {
	ID             string                          `json:"id,omitempty"`
	AgentRequestID string                          `json:"agent_request_id"`
	Type           AgentCompletionRequestEventType `json:"type"`
	LLMResponse    string                          `json:"llm_response,omitempty"`
	Metadata       interface{}                     `json:"metadata,omitempty"`
}

func (ac AgentCompletionRequestEvent) TableName() StorageTableName {
	return StorageTableNameAgentCompletionRequestEvents
}
