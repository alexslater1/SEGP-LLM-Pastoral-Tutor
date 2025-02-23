package agent

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/storage"
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

func TestNewOffloadTaskEvent(t *testing.T) {
	requestID := "req123"
	entityId := "ent123"
	task := "test task"

	event := NewOffloadTaskEvent(context_keys.SetRequestID(context.Background(), requestID), entityId, task)

	assert.Equal(t, AgentEventTypeOffloadTask, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, entityId, event.Data["entityId"])
	assert.Equal(t, task, event.Data["task"])
}

func TestNewAnswerErrorEvent(t *testing.T) {
	requestID := "req123"
	error := "test error"

	event := NewAnswerErrorEvent(context_keys.SetRequestID(context.Background(), requestID), error)

	assert.Equal(t, AgentEventTypeAnswerError, event.Type)
	assert.Equal(t, requestID, event.RequestID)
	assert.Equal(t, error, event.Data["error"])
}

func TestIdk(t *testing.T) {
	if os.Getenv("TEST_IDK") != "true" {
		t.Skip("skipping test")
	}

	storage := storage.NewMemoryStorage()
	agent := NewWellbeingMentalHealthPersonalDevelopmentAgent()
	// llm := llm.NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	a := NewLoggingAgent(NewEventStoringAgent(agent, storage))

	resp, err := a.Run(context.Background(), "test ")
	if err != nil {
		t.Fatalf("error running agent:  %v", err)
	}

	fmt.Printf("%+v\n", resp)
}
