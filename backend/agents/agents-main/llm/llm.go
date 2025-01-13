package llm

import "context"

type LLM interface {
	ChatCompletion(ctx context.Context, prompt string, model string) (string, error)
}
