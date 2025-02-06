package tools

import (
	"encoding/json"
	"fmt"
	"sort"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
)

type ToolHandler struct {
	Tools map[string]Tool
}

func NewDefaultToolHandler(googleSearch googleSearch.GoogleSearchClient, knowledge knowledge.Knowledge) *ToolHandler {
	return NewToolHandler(
		[]Tool{
			NewGoogleSearchFirstResultsPageContentsTool(googleSearch, 3),
			NewRagTool(knowledge),
			NewGoogleMapsResultsTool(googleSearch),
			NewGoogleMapsPlaceTool(googleSearch),
		},
	)
}

func NewGoogleSearchToolHandler(googleSearch googleSearch.GoogleSearchClient) *ToolHandler {
	return NewToolHandler(
		[]Tool{
			NewGoogleSearchFirstResultsPageContentsTool(googleSearch, 3),
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

	case *GoogleMapsResultsTool:
		query, ok := parsedArgs["query"]
		if !ok {
			return nil, fmt.Errorf("query is required")
		}
		return typedTool.GoogleMapsResultsFor(query.(string))

	case *GoogleMapsPlaceTool:
		googleMapsPlaceURL, ok := parsedArgs["google_maps_place_url"]
		if !ok {
			return nil, fmt.Errorf("google_maps_place_url is required")
		}
		return typedTool.PlaceDetailsFor(googleMapsPlaceURL.(string))

	case *GoogleSearchFirstResultsPageContentsTool:
		query, ok := parsedArgs["query"]
		if !ok {
			return nil, fmt.Errorf("query is required")
		}
		return typedTool.GoogleSearchFirstResultsPageContentsFor(query.(string))
	}
	return nil, fmt.Errorf("no tool matched the name %s", toolCall.Name)
}

func (t *ToolHandler) ToolDefinitions() []ToolDefinition {
	toolDefinitions := []ToolDefinition{}
	for _, tool := range t.Tools {
		toolDefinitions = append(toolDefinitions, tool.Definition())
	}

	// Sort toolDefinitions by Name
	sort.Slice(toolDefinitions, func(i, j int) bool {
		return toolDefinitions[i].Name < toolDefinitions[j].Name
	})
	
	return toolDefinitions
}

func (t *ToolHandler) ToolDefinitionsWithDescribingAction() []ToolDefinition {
	return t.toolDefinitionsWithDescribingAction()
}

func (t *ToolHandler) toolDefinitionsWithDescribingAction() []ToolDefinition {
	toolDefinitions := []ToolDefinition{}

	describingAction := Parameter{
		Name:        "descriptionOfAction",
		Description: "A short description of the action you are taking. Present continuous tense, and must be specific to the tool being used. For example \"Searching Google results for italian restaurants near Imperial College London\".",
		Type:        ParameterTypeString,
	}

	for _, tool := range t.Tools {
		definition := tool.Definition()
		newDefinition := ToolDefinition{
			Name:        definition.Name,
			Description: definition.Description,
			Parameters:  append(definition.Parameters, describingAction),
		}
		toolDefinitions = append(toolDefinitions, newDefinition)
	}
	return toolDefinitions
}

func (t *ToolHandler) ToolDefinitionsWithThoughts(toolDefinitions []ToolDefinition) []ToolDefinition {
	tds := []ToolDefinition{}

	thinkingParameter := Parameter{
		Name:        "_thinking",
		Description: "Put your thoughts here",
		Type:        ParameterTypeString,
	}

	for _, tool := range toolDefinitions {
		tds = append(tds, ToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  append([]Parameter{thinkingParameter}, tool.Parameters...),
		})
	}
	return tds
}
