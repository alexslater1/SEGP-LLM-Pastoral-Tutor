package llm

import (
	"context"

	"github.com/segp/agents-main/tools"
)

type LLM interface {
	ChatCompletion(ctx context.Context, prompt string) (*string, error)
	StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error)
	ChatCompletionWithTools(ctx context.Context, prompt string, ts []tools.Tool, toolChoice tools.ToolChoice) ([]tools.ToolCall, error)
}
