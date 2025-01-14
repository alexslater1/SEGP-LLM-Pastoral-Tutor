package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const (
	openaiModel = openai.GPT4o
)

type OpenAiLLM struct {
	client *openai.Client
}

func NewOpenAiLLM(apiKey string) *OpenAiLLM {
	return &OpenAiLLM{
		client: openai.NewClient(apiKey),
	}
}

func (o *OpenAiLLM) ChatCompletion(ctx context.Context, prompt string) (*string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openaiModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	
	if err != nil {
		return nil, fmt.Errorf("CreateChatCompletion error: %v", err)
	}

	return &resp.Choices[0].Message.Content, err
}

func (o *OpenAiLLM) StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error) {
	responseSchema, err := jsonschema.GenerateSchemaForType(schema)
	if err != nil {
		return nil, fmt.Errorf("GenerateSchemaForType error: %v", err)
	}

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openaiModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "response_format",
				Schema: responseSchema,
				Strict: true,
			},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("CreateChatCompletion error: %v", err)
	}

	return &resp.Choices[0].Message.Content, err
}
