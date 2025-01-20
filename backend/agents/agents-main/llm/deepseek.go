package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cohesion-org/deepseek-go"
	"github.com/segp/agents-main/tools"
	"github.com/segp/agents-main/utils"
)

type DeepSeekLLM struct {
	client *deepseek.Client
	apiKey string
}

func NewDeepSeekLLM(apiKey string) *DeepSeekLLM {
	utils.Required(apiKey, "DEEPSEEK_API_KEY")
	return &DeepSeekLLM{client: deepseek.NewClient(apiKey), apiKey: apiKey}
}

func (d *DeepSeekLLM) ChatCompletion(ctx context.Context, prompt string) (*string, error) {
	resp, err := d.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: prompt},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("DeepSeek chat completion error: %v", err)
	}

	return &resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error) {
	jsonschema, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling schema: %v", err)
	}

	injectedPrompt := fmt.Sprintf("%s\n\nPlease output a JSON object that matches the following schema: %s", prompt, string(jsonschema))

	resp, err := d.client.CreateChatCompletion(ctx, &deepseek.ChatCompletionRequest{
		Model: deepseek.DeepSeekChat,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleUser, Content: injectedPrompt},
		},
		ResponseFormat: &deepseek.ResponseFormat{Type: "json_object"},
	})

	if err != nil {
		return nil, fmt.Errorf("DeepSeek chat completion error: %v", err)
	}

	return &resp.Choices[0].Message.Content, nil
}

func (d *DeepSeekLLM) ChatCompletionWithTools(ctx context.Context, prompt string, tools []tools.ToolDefinition, toolChoice tools.ToolChoice) ([]tools.ToolCall, error) {
	reqBodyStr := requestBodyStrFrom(tools, deepseek.DeepSeekChat, prompt, toolChoice)

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", strings.NewReader(reqBodyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.apiKey))

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return extractToolCalls(string(body))
}
