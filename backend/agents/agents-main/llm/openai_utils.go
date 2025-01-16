package llm

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/segp/agents-main/tools"
)

type OpenaiToolPropertyType string

const (
	OpenaiToolPropertyTypeString  OpenaiToolPropertyType = "string"
	OpenaiToolPropertyTypeNumber  OpenaiToolPropertyType = "number"
	OpenaiToolPropertyTypeBoolean OpenaiToolPropertyType = "boolean"
	OpenaiToolPropertyTypeArray   OpenaiToolPropertyType = "array"
)

type OpenaiToolProperty struct {
	Type        OpenaiToolPropertyType `json:"type"`
	Description string                 `json:"description"`
}

type OpenaiToolParameters struct {
	Type                 string                        `json:"type"`
	Properties           map[string]OpenaiToolProperty `json:"properties"`
	Required             []string                      `json:"required"`
	AdditionalProperties bool                          `json:"additionalProperties"`
}

func openaiToolFrom(tool tools.ToolDefinition) openai.Tool {
	fn := openai.FunctionDefinition{
		Name:        tool.Name,
		Description: tool.Description,
		Parameters:  openaiParamsFrom(tool.Parameters),
		Strict:      true,
	}

	return openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &fn,
	}
}

func openaiToolPropertyTypeFrom(paramType tools.ParameterType) OpenaiToolPropertyType {
	switch paramType {
	case tools.ParameterTypeString:
		return OpenaiToolPropertyTypeString
	case tools.ParameterTypeNumber:
		return OpenaiToolPropertyTypeNumber
	case tools.ParameterTypeBoolean:
		return OpenaiToolPropertyTypeBoolean
	case tools.ParameterTypeArray:
		return OpenaiToolPropertyTypeArray
	}

	panic(fmt.Sprintf("Invalid parameter type: %v", paramType))
}

func openaiParamsFrom(params []tools.Parameter) OpenaiToolParameters {
	properties := map[string]OpenaiToolProperty{}
	for _, param := range params {
		properties[param.Name] = OpenaiToolProperty{
			Type:        openaiToolPropertyTypeFrom(param.Type),
			Description: param.Description,
		}
	}

	required := []string{}
	for _, param := range params {
		required = append(required, param.Name)
	}

	return OpenaiToolParameters{
		Type:                 "object",
		Properties:           properties,
		Required:             required,
		AdditionalProperties: false,
	}
}

func toolCallFromOpenai(toolCall openai.ToolCall) tools.ToolCall {
	return tools.ToolCall{
		Name:      toolCall.Function.Name,
		Arguments: toolCall.Function.Arguments,
	}
}

func openaiToolChoiceFrom(toolChoice tools.ToolChoice) any {
	switch toolChoice.Type {
	case tools.ToolChoiceTypeAuto:
		return "auto"
	case tools.ToolChoiceTypeRequired:
		return "required"
	case tools.ToolChoiceTypeForcedOne:
		return openai.ToolChoice{Type: "function", Function: openai.ToolFunction{Name: toolChoice.FunctionName}}
	}

	panic(fmt.Sprintf("Invalid tool choice type: %v", toolChoice.Type))
}
