package tools

import (
	"fmt"

	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/utils"
)

type SearchTool struct {
	knowledge               knowledge.Knowledge
	googleSearchResultsTool *GoogleSearchFirstResultsPageContentsTool
}

func NewSearchTool(knowledge knowledge.Knowledge, googleSearchResultsTool *GoogleSearchFirstResultsPageContentsTool) *SearchTool {
	return &SearchTool{
		knowledge:               knowledge,
		googleSearchResultsTool: googleSearchResultsTool,
	}
}

func (s *SearchTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "search_tool",
		Description: "Search for information about a given query. Uses a combination of RAG (on Imperial College Data) and Google Search.",
		Parameters: []Parameter{
			{Name: "query", Type: "string", Description: "The query to search for."},
		},
	}
}

func (s *SearchTool) Search(query string) (*string, error) {
	ragTask := utils.DoAsync(func() (*string, error) {
		return s.knowledge.Get(query)
	})

	googleSearchResults, err := s.googleSearchResultsTool.GoogleSearchFirstResultsPageContentsFor(query)
	if err != nil {
		return nil, err
	}

	ragResult, err := ragTask.Get()
	if err != nil {
		return nil, err
	}

	res := fmt.Sprintf("Imperial RAG Data:\n%s\n\nGoogle Search Data:\n%s", *ragResult, *googleSearchResults)
	return &res, nil
}
