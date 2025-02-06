package agent

import (
	"testing"

	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateRequestEvent(t *testing.T) {
	requestID := "req123"
	query := "test query"

	event := NewCreateRequestEvent(requestID, query)

	assert.Equal(t, AgentEventTypeCreateRequest, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, query, event.Data["query"])
}

func TestNewThinkEvent(t *testing.T) {
	requestID := "req123"
	thoughts := "thinking process"

	event := NewThinkEvent(requestID, thoughts)

	assert.Equal(t, AgentEventTypeThink, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, thoughts, event.Data["thoughts"])
}

func TestNewToolCallResultEvent(t *testing.T) {
	requestID := "req123"
	result := "tool result"

	event := NewToolCallResultEvent(requestID, &result)

	assert.Equal(t, AgentEventTypeToolCallResult, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, result, event.Data["toolCallResult"])
}

func TestNewAnswerSuccessEvent(t *testing.T) {
	requestID := "req123"
	answer := "test answer"
	reason := "test reason"

	event := NewAnswerSuccessEvent(requestID, answer, reason)

	assert.Equal(t, AgentEventTypeAnswerSuccess, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, answer, event.Data["answer"])
	assert.Equal(t, reason, event.Data["reason"])
}

func TestNewObservationEvent(t *testing.T) {
	requestID := "req123"
	observation := "test observation"

	event := NewObservationEvent(requestID, observation)

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

	event := NewToolCallChoiceEvent(requestID, toolCall)

	assert.Equal(t, AgentEventTypeToolCallChoice, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, toolCall, event.Data["toolCallChoice"])
}
