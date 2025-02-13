package agent

import (
	"context"
	"testing"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryEvent(t *testing.T) {
	requestID := "req123"
	query := "test query"

	event := NewQueryEvent(context_keys.SetRequestID(context.Background(), requestID), query)

	assert.Equal(t, AgentEventTypeQuery, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, query, event.Data["query"])
}

func TestNewAnswerErrorEvent(t *testing.T) {
	requestID := "req123"
	error := "test error"

	event := NewAnswerErrorEvent(context_keys.SetRequestID(context.Background(), requestID), error)

	assert.Equal(t, AgentEventTypeAnswerError, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, error, event.Data["error"])
}

func TestNewThinkEvent(t *testing.T) {
	requestID := "req123"
	thoughts := "thinking process"

	event := NewThinkEvent(context_keys.SetRequestID(context.Background(), requestID), thoughts)

	assert.Equal(t, AgentEventTypeThink, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, thoughts, event.Data["thoughts"])
}

func TestNewToolCallResultEvent(t *testing.T) {
	requestID := "req123"
	result := "tool result"

	event := NewToolCallResultEvent(context_keys.SetRequestID(context.Background(), requestID), &result)

	assert.Equal(t, AgentEventTypeToolCallResult, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, result, event.Data["toolCallResult"])
}

func TestNewAnswerSuccessEvent(t *testing.T) {
	requestID := "req123"
	answer := "test answer"
	reason := "test reason"

	event := NewAnswerSuccessEvent(context_keys.SetRequestID(context.Background(), requestID), answer, reason)

	assert.Equal(t, AgentEventTypeAnswerSuccess, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, answer, event.Data["answer"])
	assert.Equal(t, reason, event.Data["reason"])
}

func TestNewObservationEvent(t *testing.T) {
	requestID := "req123"
	observation := "test observation"

	event := NewObservationEvent(context_keys.SetRequestID(context.Background(), requestID), observation)

	assert.Equal(t, AgentEventTypeObservation, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, observation, event.Data["observation"])
}

func TestNewToolCallChoiceEvent(t *testing.T) {
	requestID := "req123"
	toolCall := tools.ToolCall{
		Name:      "testTool",
		Arguments: "test arguments",
	}

	event := NewToolCallChoiceEvent(context_keys.SetRequestID(context.Background(), requestID), toolCall)

	assert.Equal(t, AgentEventTypeToolCallChoice, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, toolCall, event.Data["toolCallChoice"])
}
