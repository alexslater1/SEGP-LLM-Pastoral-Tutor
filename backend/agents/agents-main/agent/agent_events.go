package agent

import (
	"context"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/tools"
)

type AgentEventType string

const (
	AgentEventTypeThink           AgentEventType = "think"
	AgentEventTypeQuery           AgentEventType = "query"
	AgentEventTypeToolCallChoice  AgentEventType = "tool_call_choice"
	AgentEventTypeToolCallResult  AgentEventType = "tool_call_result"
	AgentEventTypeAnswerSuccess   AgentEventType = "answer_success"
	AgentEventTypeAnswerError     AgentEventType = "answer_error"
	AgentEventTypeOffloadTask     AgentEventType = "offload_task"
	AgentEventTypeRouterSelection AgentEventType = "router_selection"
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
	data := map[string]interface{}{"query": query}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}

	return AgentEvent{
		Type:      AgentEventTypeQuery,
		RequestID: requestId,
		Data:      data,
	}
}

func NewThinkEvent(ctx context.Context, thoughts string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"thoughts": thoughts}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeThink,
		RequestID: requestId,
		Data:      data,
	}
}

func NewToolCallResultEvent(ctx context.Context, toolCallResult *string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"toolCallResult": *toolCallResult}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeToolCallResult,
		RequestID: requestId,
		Data:      data,
	}
}

func NewAnswerSuccessEvent(ctx context.Context, answer string, reason string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"answer": answer, "reason": reason}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeAnswerSuccess,
		RequestID: requestId,
		Data:      data,
	}
}

func NewToolCallChoiceEvent(ctx context.Context, toolCall tools.ToolCall) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"toolCallChoice": toolCall}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeToolCallChoice,
		RequestID: requestId,
		Data:      data,
	}
}

func NewOffloadTaskEvent(ctx context.Context, entityId string, task string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"entityId": entityId, "task": task}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeOffloadTask,
		RequestID: requestId,
		Data:      data,
	}
}

func NewAnswerErrorEvent(ctx context.Context, error string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"error": error}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeAnswerError,
		RequestID: requestId,
		Data:      data,
	}
}

func NewRouterSelectionEvent(ctx context.Context, selectedAgentID string, reason string) AgentEvent {
	requestId, ok := context_keys.GetRequestID(ctx)
	if !ok {
		requestId = ""
	}
	data := map[string]interface{}{"selectedAgentID": selectedAgentID, "reason": reason}
	if agentID, ok := context_keys.GetAgentID(ctx); ok {
		data["agentID"] = agentID
	}
	return AgentEvent{
		Type:      AgentEventTypeRouterSelection,
		RequestID: requestId,
		Data:      data,
	}
}
