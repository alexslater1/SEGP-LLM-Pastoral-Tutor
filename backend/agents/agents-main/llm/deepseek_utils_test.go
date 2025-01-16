package llm

import (
	"testing"

	"github.com/segp/agents-main/tools"
)

func TestDeepseekToolFrom(t *testing.T) {
	prompt := "What is the weather in Paris on 2025-01-14?"
	ts := []tools.ToolDefinition{tools.CheckWeatherTool()}
	model := "deepseek-chat"

	requestBodyStr := requestBodyStrFrom(ts, model, prompt, tools.ToolChoice{Type: tools.ToolChoiceTypeForcedOne, FunctionName: "check_weather"})

	t.Logf("RequestBodyStr: %v", requestBodyStr)
}
