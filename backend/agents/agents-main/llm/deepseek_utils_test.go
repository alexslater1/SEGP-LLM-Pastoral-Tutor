package llm

import (
	"testing"

	"github.com/segp/agents-main/tools"
)

func TestDeepseekToolFrom(t *testing.T) {
	prompt := "What is the weather in Paris on 2025-01-14?"
	tools := []tools.Tool{tools.CheckWeatherTool()}
	model := "deepseek-chat"

	requestBodyStr := requestBodyStrFrom(tools, model, prompt)

	t.Logf("RequestBodyStr: %v", requestBodyStr)
}
