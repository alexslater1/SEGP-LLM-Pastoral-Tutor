package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cohesion-org/deepseek-go"
)

type DeepSeekLLM struct {
	client *deepseek.Client
}

func NewDeepSeekLLM(apiKey string) *DeepSeekLLM {
	return &DeepSeekLLM{client: deepseek.NewClient(apiKey)}
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
