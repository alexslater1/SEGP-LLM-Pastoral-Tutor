package llm

import (
	"context"
	"testing"

	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestMockLLM_ChatCompletion(t *testing.T) {
	tests := []struct {
		name          string
		setupChain    func(*MockLLM) *MockLLM
		expectedResp  string
		expectPanic   bool
		panicMessage  string
	}{
		{
			name: "successful chat completion",
			setupChain: func(m *MockLLM) *MockLLM {
				return m.NewCallChain().ThenChat("Hello, world!").Set()
			},
			expectedResp: "Hello, world!",
		},
		{
			name: "wrong response type",
			setupChain: func(m *MockLLM) *MockLLM {
				return m.NewCallChain().ThenStructured("wrong type").Set()
			},
			expectPanic:  true,
			panicMessage: "expected ChatResponse, got StructuredResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := NewMockLLM()
			llm = tt.setupChain(llm)

			if tt.expectPanic {
				assert.PanicsWithValue(t, tt.panicMessage, func() {
					llm.ChatCompletion(context.Background(), "test prompt")
				})
				return
			}

			resp, err := llm.ChatCompletion(context.Background(), "test prompt")
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, *resp)
		})
	}
}

func TestMockLLM_StructuredOutputCompletion(t *testing.T) {
	tests := []struct {
		name          string
		setupChain    func(*MockLLM) *MockLLM
		expectedResp  string
		expectPanic   bool
		panicMessage  string
	}{
		{
			name: "successful structured completion",
			setupChain: func(m *MockLLM) *MockLLM {
				return m.NewCallChain().ThenStructured(`{"key": "value"}`).Set()
			},
			expectedResp: `{"key": "value"}`,
		},
		{
			name: "wrong response type",
			setupChain: func(m *MockLLM) *MockLLM {
				return m.NewCallChain().ThenChat("wrong type").Set()
			},
			expectPanic:  true,
			panicMessage: "expected StructuredResponse, got ChatResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := NewMockLLM()
			llm = tt.setupChain(llm)

			if tt.expectPanic {
				assert.PanicsWithValue(t, tt.panicMessage, func() {
					llm.StructuredOutputCompletion(context.Background(), "test prompt", nil)
				})
				return
			}

			resp, err := llm.StructuredOutputCompletion(context.Background(), "test prompt", nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResp, *resp)
		})
	}
}

func TestMockLLM_ChatCompletionWithTools(t *testing.T) {
	tests := []struct {
		name          string
		setupChain    func(*MockLLM) *MockLLM
		expectedCalls []tools.ToolCall
		expectPanic   bool
		panicMessage  string
	}{
		{
			name: "successful tool completion",
			setupChain: func(m *MockLLM) *MockLLM {
				toolCalls := []tools.ToolCall{
					{
						Name:      "test_tool",
						Arguments: `{"arg": "value"}`,
					},
				}
				return m.NewCallChain().ThenTool(toolCalls).Set()
			},
			expectedCalls: []tools.ToolCall{
				{
					Name:      "test_tool",
					Arguments: `{"arg": "value"}`,
				},
			},
		},
		{
			name: "wrong response type",
			setupChain: func(m *MockLLM) *MockLLM {
				return m.NewCallChain().ThenChat("wrong type").Set()
			},
			expectPanic:  true,
			panicMessage: "expected ToolResponse, got ChatResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llm := NewMockLLM()
			llm = tt.setupChain(llm)

			if tt.expectPanic {
				assert.PanicsWithValue(t, tt.panicMessage, func() {
					llm.ChatCompletionWithTools(context.Background(), "test prompt", nil, tools.ToolChoice{Type: tools.ToolChoiceTypeAuto})
				})
				return
			}

			resp, err := llm.ChatCompletionWithTools(context.Background(), "test prompt", nil, tools.ToolChoice{Type: tools.ToolChoiceTypeAuto})
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCalls, resp)
		})
	}
}

func TestMockLLM_ChainedResponses(t *testing.T) {
	llm := NewMockLLM()
	llm.NewCallChain().
		ThenChat("chat response").
		ThenStructured(`{"structured": true}`).
		ThenTool([]tools.ToolCall{{Name: "test", Arguments: `{"arg": "value"}`}}).
		Set()

	// Test chat response
	chatResp, err := llm.ChatCompletion(context.Background(), "prompt")
	assert.NoError(t, err)
	assert.Equal(t, "chat response", *chatResp)

	// Test structured response
	structResp, err := llm.StructuredOutputCompletion(context.Background(), "prompt", nil)
	assert.NoError(t, err)
	assert.Equal(t, `{"structured": true}`, *structResp)

	// Test tool response
	toolResp, err := llm.ChatCompletionWithTools(context.Background(), "prompt", nil, tools.ToolChoice{Type: tools.ToolChoiceTypeAuto})
	assert.NoError(t, err)
	assert.Equal(t, []tools.ToolCall{{Name: "test", Arguments: `{"arg": "value"}`}}, toolResp)
}
