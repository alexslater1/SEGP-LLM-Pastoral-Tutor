package tools

import (
	"encoding/json"
	"fmt"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
)

type ToolHandler struct {
	Tools map[string]Tool
}

func NewDefaultToolHandler(googleSearch googleSearch.GoogleSearchClient, knowledge knowledge.Knowledge) *ToolHandler {
	return NewToolHandler(
		[]Tool{
			NewGoogleSearchUrlTool(googleSearch),
			NewGoogleSearchResultsTool(googleSearch),
			NewRagTool(knowledge),
		},
	)
}

func NewGoogleSearchToolHandler(googleSearch googleSearch.GoogleSearchClient) *ToolHandler {
	return NewToolHandler(
		[]Tool{
			NewGoogleSearchUrlTool(googleSearch),
			NewGoogleSearchResultsTool(googleSearch),
		},
	)
}

func NewToolHandler(tools []Tool) *ToolHandler {
	toolMap := make(map[string]Tool)
	for _, tool := range tools {
		toolMap[tool.Definition().Name] = tool
	}

	if _, ok := toolMap["no_tool"]; !ok {
		toolMap["no_tool"] = NewNoTool()
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
		query, ok := parsedArgs["query"]
		if !ok {
			return nil, fmt.Errorf("query is required")
		}
		return typedTool.GoogleSearchResultsFor(query.(string))
	case *GoogleSearchUrlTool:
		url, ok := parsedArgs["URL"]
		if !ok {
			return nil, fmt.Errorf("URL is required")
		}
		return typedTool.PageContentFor(url.(string))
	case *DateTool:
		date := typedTool.GetCurrentDate()
		return &date, nil
	case *RagTool:
		query, ok := parsedArgs["query"]
		if !ok {
			return nil, fmt.Errorf("query is required")
		}
		return typedTool.SearchRagFor(query.(string))
	}

	return nil, fmt.Errorf("no tool matched the name %s", toolCall.Name)
}

func (t *ToolHandler) ToolDefinitions() []ToolDefinition {
	toolDefinitions := []ToolDefinition{}
	for _, tool := range t.Tools {
		toolDefinitions = append(toolDefinitions, tool.Definition())
	}
	return toolDefinitions
}
