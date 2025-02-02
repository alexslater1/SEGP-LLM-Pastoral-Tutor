package llm

import (
	"context"
	"testing"

	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestMockLLM(t *testing.T) {
	t.Run("ChatCompletion", func(t *testing.T) {
		mock := NewMockLLM()
		expectedResponse := "Hello, world!"
		mock.SetChatResponse("test prompt", expectedResponse)

		// Test successful case
		response, err := mock.ChatCompletion(context.Background(), "test prompt")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, *response)

		// Test error case for unknown prompt
		response, err = mock.ChatCompletion(context.Background(), "unknown prompt")
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "no mock response set for prompt")
	})

	t.Run("StructuredOutputCompletion", func(t *testing.T) {
		mock := NewMockLLM()
		expectedResponse := `{"key": "value"}`
		mock.SetStructuredResponse("test prompt", expectedResponse)

		// Test successful case
		response, err := mock.StructuredOutputCompletion(context.Background(), "test prompt", nil)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedResponse, *response)

		// Test error case for unknown prompt
		response, err = mock.StructuredOutputCompletion(context.Background(), "unknown prompt", nil)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "no mock response set for prompt")
	})

	t.Run("ChatCompletionWithTools", func(t *testing.T) {
		mock := NewMockLLM()
		expectedToolCalls := []tools.ToolCall{
			{
				Name:      "calculator",
				Arguments: "{\"x\": 1, \"y\": 2}",
			},
		}
		mock.SetToolResponse("test prompt", expectedToolCalls)

		// Test successful case
		toolDefs := []tools.ToolDefinition{{
			Name:        "calculator",
			Description: "Adds two numbers",
		}}
		response, err := mock.ChatCompletionWithTools(context.Background(), "test prompt", toolDefs, tools.ToolChoice{})
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedToolCalls, response)

		// Test error case for unknown prompt
		response, err = mock.ChatCompletionWithTools(context.Background(), "unknown prompt", toolDefs, tools.ToolChoice{})
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "no mock response set for prompt")
	})
}
