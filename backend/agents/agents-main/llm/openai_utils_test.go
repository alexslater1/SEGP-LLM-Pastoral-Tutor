package llm

import (
	"testing"

	"github.com/sashabaranov/go-openai"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestOpenaiToolPropertyTypeFrom(t *testing.T) {
	// Define a weather tool with coordinates parameters
	tool := tools.Tool{
		Name:        "get_weather",
		Description: "Get current temperature for provided coordinates in celsius.",
		Parameters: []tools.Parameter{
			{
				Name:        "latitude",
				Description: "Latitude coordinate",
				Type:        tools.ParameterTypeNumber,
			},
			{
				Name:        "longitude",
				Description: "Longitude coordinate",
				Type:        tools.ParameterTypeNumber,
			},
		},
	}

	openaiTool := openaiToolFrom(tool)

	// Verify the basic tool properties
	assert.Equal(t, openaiTool, openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_weather",
			Description: "Get current temperature for provided coordinates in celsius.",
			Parameters: OpenaiToolParameters{
				Type: "object",
				Properties: map[string]OpenaiToolProperty{
					"latitude":  {Type: OpenaiToolPropertyTypeNumber},
					"longitude": {Type: OpenaiToolPropertyTypeNumber},
				},
				Required:             []string{"latitude", "longitude"},
				AdditionalProperties: false,
			},
			Strict: true,
		},
	})
}

func TestOpenaiToolChoiceFrom(t *testing.T) {
	toolChoice := tools.ToolChoice{Type: tools.ToolChoiceTypeForcedOne, FunctionName: "get_weather"}
	openaiToolChoice := openaiToolChoiceFrom(toolChoice)
	assert.Equal(t, openaiToolChoice, openai.ToolChoice{Type: "function", Function: openai.ToolFunction{Name: "get_weather"}})

	toolChoice = tools.ToolChoice{Type: tools.ToolChoiceTypeRequired}
	openaiToolChoice = openaiToolChoiceFrom(toolChoice)
	assert.Equal(t, openaiToolChoice, "required")

	toolChoice = tools.ToolChoice{Type: tools.ToolChoiceTypeAuto}
	openaiToolChoice = openaiToolChoiceFrom(toolChoice)
	assert.Equal(t, openaiToolChoice, "auto")
}
