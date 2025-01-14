package llm

import (
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/segp/agents-main/tools"
	"github.com/stretchr/testify/assert"
)

func TestGenaiToolPropertyTypeFrom(t *testing.T) {
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

	genaiTool := genaiToolFrom(tool)

	// Verify the basic tool properties
	assert.Equal(t, *genaiTool, genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_weather",
				Description: "Get current temperature for provided coordinates in celsius.",
				Parameters:  &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"latitude":  {Type: genai.TypeNumber},
						"longitude": {Type: genai.TypeNumber},
					},
					Required: []string{"latitude", "longitude"},
				},
			},
		},
	})
}

func TestFunctionCallingConfig(t *testing.T) {
	toolChoice := tools.ToolChoice{Type: tools.ToolChoiceTypeForcedOne, FunctionName: "get_weather"}
	functionCallingConfig := functionCallingConfigFrom(toolChoice)
	assert.Equal(t, *functionCallingConfig, genai.FunctionCallingConfig{Mode: genai.FunctionCallingAny, AllowedFunctionNames: []string{"get_weather"}})

	toolChoice = tools.ToolChoice{Type: tools.ToolChoiceTypeRequired}
	functionCallingConfig = functionCallingConfigFrom(toolChoice)
	assert.Equal(t, *functionCallingConfig, genai.FunctionCallingConfig{Mode: genai.FunctionCallingAny})

	toolChoice = tools.ToolChoice{Type: tools.ToolChoiceTypeAuto}
	functionCallingConfig = functionCallingConfigFrom(toolChoice)
	assert.Equal(t, *functionCallingConfig, genai.FunctionCallingConfig{Mode: genai.FunctionCallingAuto})
}