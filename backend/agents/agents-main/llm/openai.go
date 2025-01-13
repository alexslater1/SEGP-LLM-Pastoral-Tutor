package llm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sashabaranov/go-openai"
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

func (o *OpenAiLLM) ChatCompletion(ctx context.Context, prompt string) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openaiModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	return resp.Choices[0].Message.Content, err
}

func (o *OpenAiLLM) StructuredChatCompletion(ctx context.Context, prompt string, model string) (string, error) {
	return "", nil
}

func iterateJSON(data interface{}, indent string) {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		// Iterate over map keys
		for key, value := range data.(map[string]interface{}) {
			fmt.Printf("%sKey: %s\n", indent, key)
			iterateJSON(value, indent+"  ")
		}
	case reflect.Slice:
		// Iterate over array/slice
		for i, value := range data.([]interface{}) {
			fmt.Printf("%sIndex %d:\n", indent, i)
			iterateJSON(value, indent+"  ")
		}
	default:
		// Print primitive values
		fmt.Printf("%sValue: %v\n", indent, data)
	}
}
