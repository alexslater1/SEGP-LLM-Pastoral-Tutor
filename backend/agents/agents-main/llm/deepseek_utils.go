package llm

import (
	"encoding/json"
	"fmt"

	"github.com/segp/agents-main/tools"
)

type PropertyType string

const (
	PropertyTypeString  PropertyType = "string"
	PropertyTypeNumber  PropertyType = "number"
	PropertyTypeBoolean PropertyType = "boolean"
	PropertyTypeArray   PropertyType = "array"
	PropertyTypeObject  PropertyType = "object"
)

type ParameterType string

const (
	ParameterTypeObject ParameterType = "object"
)

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type Property struct {
	Type        PropertyType `json:"type"`
	Description string       `json:"description"`
}

type Parameters struct {
	Type       ParameterType       `json:"type"`
	Required   []string            `json:"required"`
	Properties map[string]Property `json:"properties"`
}

type DeepSeekFunction struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}

type DeepSeekTool struct {
	Type     ToolType         `json:"type"`
	Function DeepSeekFunction `json:"function"`
}

func deepSeekToolFrom(tool tools.ToolDefinition) DeepSeekTool {

	required := []string{}
	properties := map[string]Property{}
	for _, param := range tool.Parameters {
		required = append(required, param.Name)
		properties[param.Name] = Property{
			Type:        PropertyTypeString,
			Description: param.Description,
		}
	}

	return DeepSeekTool{
		Type: ToolTypeFunction,
		Function: DeepSeekFunction{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters: Parameters{
				Type:       ParameterTypeObject,
				Required:   required,
				Properties: properties,
			},
		},
	}
}

func toolsChoiceStrFrom(toolChoice tools.ToolChoice) any {
	switch toolChoice.Type {
	case tools.ToolChoiceTypeAuto:
		return fmt.Sprintf("auto")
	case tools.ToolChoiceTypeRequired:
		return fmt.Sprintf("required")
	case tools.ToolChoiceTypeForcedOne:
		return map[string]interface{}{
			"type": "function",
			"function": map[string]string{
				"name": toolChoice.FunctionName,
			},
		}
	}

	panic(fmt.Sprintf("unknown tool choice type: %v", toolChoice.Type))
}

func requestBodyStrFrom(tools []tools.ToolDefinition, model string, prompt string, toolChoice tools.ToolChoice) string {
	deepSeekTools := []DeepSeekTool{}
	for _, tool := range tools {
		deepSeekTools = append(deepSeekTools, deepSeekToolFrom(tool))
	}

	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{{
			"role":    "user",
			"content": prompt,
		},
		},
		"tools":       deepSeekTools,
		"tool_choice": toolsChoiceStrFrom(toolChoice),
	}

	json, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	return string(json)
}

// ExtractToolCalls extracts tool calls from a JSON response.
func extractToolCalls(response string) ([]tools.ToolCall, error) {
	var parsedResponse struct {
		Choices []struct {
			Message struct {
				ToolCalls []struct {
					Function tools.ToolCall `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}

	err := json.Unmarshal([]byte(response), &parsedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	var toolCalls []tools.ToolCall
	for _, choice := range parsedResponse.Choices {
		for _, toolCall := range choice.Message.ToolCalls {
			toolCalls = append(toolCalls, toolCall.Function)
		}
	}

	return toolCalls, nil
}
