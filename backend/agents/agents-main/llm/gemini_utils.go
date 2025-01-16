package llm

import (
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/segp/agents-main/tools"
)

type GenaiToolPropertyType string

func genaiToolFrom(tool tools.ToolDefinition) *genai.Tool {
	fn := genai.FunctionDeclaration{
		Name:        tool.Name,
		Description: tool.Description,
		Parameters:  genaiParamsFrom(tool.Parameters),
	}

	return &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{&fn},
	}
}

func genaiToolPropertyTypeFrom(paramType tools.ParameterType) genai.Type {
	switch paramType {
	case tools.ParameterTypeString:
		return genai.TypeString
	case tools.ParameterTypeNumber:
		return genai.TypeNumber
	case tools.ParameterTypeBoolean:
		return genai.TypeBoolean
	case tools.ParameterTypeArray:
		return genai.TypeArray
	}

	panic(fmt.Sprintf("Invalid parameter type: %v", paramType))
}

func genaiParamsFrom(params []tools.Parameter) *genai.Schema {
	properties := map[string]*genai.Schema{}
	for _, param := range params {
		properties[param.Name] = &genai.Schema{
			Type:        genaiToolPropertyTypeFrom(param.Type),
			Description: param.Description,
		}
	}

	required := []string{}
	for _, param := range params {
		required = append(required, param.Name)
	}

	return &genai.Schema{
		Type:       genai.TypeObject,
		Properties: properties,
		Required:   required,
	}
}

func toolCallFromGenai(toolCall genai.FunctionCall) tools.ToolCall {

	jsonBytes, err := json.Marshal(toolCall.Args)
	if err != nil {
		panic(fmt.Sprintf("Error converting map to JSON: %v", err))

	}

	return tools.ToolCall{
		Name:      toolCall.Name,
		Arguments: string(jsonBytes),
	}
}

func functionCallingConfigFrom(toolChoice tools.ToolChoice) *genai.FunctionCallingConfig {
	switch toolChoice.Type {
	case tools.ToolChoiceTypeAuto:
		return &genai.FunctionCallingConfig{Mode: genai.FunctionCallingAuto}
	case tools.ToolChoiceTypeRequired:
		return &genai.FunctionCallingConfig{Mode: genai.FunctionCallingAny}
	case tools.ToolChoiceTypeForcedOne:
		return &genai.FunctionCallingConfig{Mode: genai.FunctionCallingAny, AllowedFunctionNames: []string{toolChoice.FunctionName}}
	}

	panic(fmt.Sprintf("Invalid tool choice type: %v", toolChoice.Type))
}
