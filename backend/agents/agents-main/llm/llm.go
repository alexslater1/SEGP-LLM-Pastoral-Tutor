package llm

import "context"

type LLM interface {
	ChatCompletion(ctx context.Context, prompt string) (*string, error)
	StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error)
}
