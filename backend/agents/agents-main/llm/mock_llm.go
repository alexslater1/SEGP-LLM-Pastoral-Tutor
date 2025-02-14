package llm

import (
	"context"
	"fmt"

	"github.com/segp/agents-main/tools"
)

type MockLLM struct {
	// Maps to store expected responses for different inputs
	responses []LLMResponse
}

type MockLLMResponseType string

const (
	ChatResponse       MockLLMResponseType = "ChatResponse"
	StructuredResponse MockLLMResponseType = "StructuredResponse"
	ToolResponse       MockLLMResponseType = "ToolResponse"
)

type LLMResponse struct {
	Type     MockLLMResponseType
	Response any
}

type MockLLMCallChainBuilder struct {
	llm       *MockLLM
	responses []LLMResponse
}

// NewMockLLM creates a new MockLLM instance
func NewMockLLM() *MockLLM {
	return &MockLLM{
		responses: make([]LLMResponse, 0),
	}
}

func (m *MockLLM) NewCallChain() *MockLLMCallChainBuilder {
	return &MockLLMCallChainBuilder{
		llm:       m,
		responses: make([]LLMResponse, 0),
	}
}

func (m *MockLLMCallChainBuilder) ThenChat(response string) *MockLLMCallChainBuilder {
	m.responses = append(m.responses, LLMResponse{Type: ChatResponse, Response: response})
	return m
}

func (m *MockLLMCallChainBuilder) ThenStructured(response string) *MockLLMCallChainBuilder {
	m.responses = append(m.responses, LLMResponse{Type: StructuredResponse, Response: response})
	return m
}

func (m *MockLLMCallChainBuilder) ThenTool(response []tools.ToolCall) *MockLLMCallChainBuilder {
	m.responses = append(m.responses, LLMResponse{Type: ToolResponse, Response: response})
	return m
}

func (m *MockLLMCallChainBuilder) Set() *MockLLM {
	m.llm.responses = append(m.llm.responses, m.responses...)
	return m.llm
}

// ChatCompletion implements the LLM interface
func (m *MockLLM) ChatCompletion(ctx context.Context, prompt string) (*string, error) {
	response := m.responses[0]
	m.responses = m.responses[1:]

	if response.Type != ChatResponse {
		panic(fmt.Sprintf("expected ChatResponse, got %s", response.Type))
	}

	str, ok := response.Response.(string)
	if !ok {
		panic(fmt.Sprintf("expected string, got %T", response.Response))
	}

	return &str, nil
}

// StructuredOutputCompletion implements the LLM interface
func (m *MockLLM) StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error) {

	response := m.responses[0]
	m.responses = m.responses[1:]

	if response.Type != StructuredResponse {
		panic(fmt.Sprintf("expected StructuredResponse, got %s", response.Type))
	}

	str, ok := response.Response.(string)
	if !ok {
		panic(fmt.Sprintf("expected string, got %T", response.Response))
	}

	return &str, nil
}

// ChatCompletionWithTools implements the LLM interface
func (m *MockLLM) ChatCompletionWithTools(ctx context.Context, prompt string, ts []tools.ToolDefinition, toolChoice tools.ToolChoice) ([]tools.ToolCall, error) {
	response := m.responses[0]
	m.responses = m.responses[1:]

	if response.Type != ToolResponse {
		panic(fmt.Sprintf("expected ToolResponse, got %s", response.Type))
	}

	toolCalls, ok := response.Response.([]tools.ToolCall)
	if !ok {
		panic(fmt.Sprintf("expected []tools.ToolCall, got %T", response.Response))
	}

	return toolCalls, nil
}
