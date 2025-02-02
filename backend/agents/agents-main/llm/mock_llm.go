package llm

import (
	"context"
	"fmt"

	"github.com/segp/agents-main/tools"
)

type MockLLM struct {
	// Maps to store expected responses for different inputs
	chatResponses       map[string]*string
	structuredResponses map[string]*string
	toolResponses       map[string][]tools.ToolCall
}

// NewMockLLM creates a new MockLLM instance
func NewMockLLM() *MockLLM {
	return &MockLLM{
		chatResponses:       make(map[string]*string),
		structuredResponses: make(map[string]*string),
		toolResponses:       make(map[string][]tools.ToolCall),
	}
}

// SetChatResponse sets the expected response for a given chat prompt
func (m *MockLLM) SetChatResponse(prompt string, response string) {
	m.chatResponses[prompt] = &response
}

// SetStructuredResponse sets the expected response for a structured output prompt
func (m *MockLLM) SetStructuredResponse(prompt string, response string) {
	m.structuredResponses[prompt] = &response
}

// SetToolResponse sets the expected tool calls for a given prompt
func (m *MockLLM) SetToolResponse(prompt string, response []tools.ToolCall) {
	m.toolResponses[prompt] = response
}

// ChatCompletion implements the LLM interface
func (m *MockLLM) ChatCompletion(ctx context.Context, prompt string) (*string, error) {
	if response, ok := m.chatResponses[prompt]; ok {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response set for prompt: %s", prompt)
}

// StructuredOutputCompletion implements the LLM interface
func (m *MockLLM) StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error) {
	if response, ok := m.structuredResponses[prompt]; ok {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response set for prompt: %s", prompt)
}

// ChatCompletionWithTools implements the LLM interface
func (m *MockLLM) ChatCompletionWithTools(ctx context.Context, prompt string, ts []tools.ToolDefinition, toolChoice tools.ToolChoice) ([]tools.ToolCall, error) {
	if response, ok := m.toolResponses[prompt]; ok {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response set for prompt: %s", prompt)
}
