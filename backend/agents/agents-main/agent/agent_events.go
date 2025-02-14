package agent

import (
	"context"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/tools"
)

type AgentEventType string

const (
	AgentEventTypeThink          AgentEventType = "think"
	AgentEventTypeQuery          AgentEventType = "query"
	AgentEventTypeToolCallChoice AgentEventType = "tool_call_choice"
	AgentEventTypeToolCallResult AgentEventType = "tool_call_result"
	AgentEventTypeAnswerSuccess  AgentEventType = "answer_success"
	AgentEventTypeAnswerError    AgentEventType = "answer_error"
	AgentEventTypeOffloadTask    AgentEventType = "offload_task"
)

type AgentEvent struct {
	Type      AgentEventType
	RequestID string
	Data      map[string]interface{}
}

func NewQueryEvent(ctx context.Context, query string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeQuery,
		RequestID: requestId,
		Data:      map[string]interface{}{"query": query},
	}
}

func NewThinkEvent(ctx context.Context, thoughts string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeThink,
		RequestID: requestId,
		Data: map[string]interface{}{
			"thoughts": thoughts,
		},
	}
}

func NewToolCallResultEvent(ctx context.Context, toolCallResult *string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeToolCallResult,
		RequestID: requestId,
		Data: map[string]interface{}{
			"toolCallResult": *toolCallResult,
		},
	}
}

func NewAnswerSuccessEvent(ctx context.Context, answer string, reason string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeAnswerSuccess,
		RequestID: requestId,
		Data: map[string]interface{}{
			"answer": answer,
			"reason": reason,
		},
	}
}

func NewToolCallChoiceEvent(ctx context.Context, toolCall tools.ToolCall) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeToolCallChoice,
		RequestID: requestId,
		Data: map[string]interface{}{
			"toolCallChoice": toolCall,
		},
	}
}

func NewOffloadTaskEvent(ctx context.Context, entityId string, task string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeOffloadTask,
		RequestID: requestId,
		Data:      map[string]interface{}{"entityId": entityId, "task": task},
	}
}

func NewAnswerErrorEvent(ctx context.Context, error string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	return AgentEvent{
		Type:      AgentEventTypeAnswerError,
		RequestID: requestId,
		Data:      map[string]interface{}{"error": error},
	}
}
