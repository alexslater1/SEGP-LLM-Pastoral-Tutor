package tools

import "encoding/json"

type ToolHandler struct {
	Tools map[string]Tool
}

func NewToolHandler(tools []Tool) *ToolHandler {
	toolMap := make(map[string]Tool)
	for _, tool := range tools {
		toolMap[tool.Definition().Name] = tool
	}

	return &ToolHandler{
		Tools: toolMap,
	}
}

func (t *ToolHandler) Call(toolCall ToolCall) (*string, error) {
	tool := t.Tools[toolCall.Name]

	var parsedArgs map[string]interface{}
	err := json.Unmarshal([]byte(toolCall.Arguments), &parsedArgs)
	if err != nil {
		return nil, err
	}

	switch typedTool := tool.(type) {
	case *GoogleSearchResultsTool:
		return typedTool.GoogleSearchResultsFor(parsedArgs["query"].(string))
	case *GoogleSearchUrlTool:
		return typedTool.PageContentFor(parsedArgs["url"].(string))
	}

	return nil, nil
}

func (t *ToolHandler) ToolDefinitions() []ToolDefinition {
	toolDefinitions := []ToolDefinition{}
	for _, tool := range t.Tools {
		toolDefinitions = append(toolDefinitions, tool.Definition())
	}
	return toolDefinitions
}
