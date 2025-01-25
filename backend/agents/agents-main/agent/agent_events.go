package agent

import (
	"github.com/segp/agents-main/tools"
)

type AgentEventType string

const (
	AgentEventTypeCreateRequest  AgentEventType = "create_request"
	AgentEventTypeThink          AgentEventType = "think"
	AgentEventTypeToolCallChoice AgentEventType = "tool_call_choice"
	AgentEventTypeToolCallResult AgentEventType = "tool_call_result"
	AgentEventTypeAnswerSuccess  AgentEventType = "answer_success"
	AgentEventTypeAnswerError    AgentEventType = "answer_error"
	AgentEventTypeObservation    AgentEventType = "observation"
)

type AgentEvent struct {
	Type      AgentEventType
	RequestID string
	Data      map[string]interface{}
}

func NewCreateRequestEvent(requestId string, query string) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeCreateRequest,
		RequestID: requestId,
		Data: map[string]interface{}{
			"query": query,
		},
	}
}

func NewThinkEvent(requestId string, thoughts string) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeThink,
		RequestID: requestId,
		Data: map[string]interface{}{
			"thoughts": thoughts,
		},
	}
}

func NewToolCallResultEvent(requestId string, toolCallResult *string) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeToolCallResult,
		RequestID: requestId,
		Data: map[string]interface{}{
			"toolCallResult": toolCallResult,
		},
	}
}

func NewAnswerSuccessEvent(requestId string, answer string, reason string) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeAnswerSuccess,
		RequestID: requestId,
		Data: map[string]interface{}{
			"answer": answer,
			"reason": reason,
		},
	}
}

func NewObservationEvent(requestId string, observation string) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeObservation,
		RequestID: requestId,
		Data: map[string]interface{}{
			"observation": observation,
		},
	}
}

func NewToolCallChoiceEvent(requestId string, toolCall tools.ToolCall) AgentEvent {
	return AgentEvent{
		Type:      AgentEventTypeToolCallChoice,
		RequestID: requestId,
		Data: map[string]interface{}{
			"toolCallChoice": toolCall,
		},
	}
}
