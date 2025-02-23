package tools

import (
	"fmt"
	"os"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
)

func TestSearchTool(t *testing.T) {
	ragKnowledge := knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
	googleSearchResultsTool := NewGoogleSearchFirstResultsPageContentsTool(googleSearch.NewNonHeadlessRodClient(), 3)
	searchTool := NewSearchTool(ragKnowledge, googleSearchResultsTool)

	result, err := searchTool.Search("who is the imperial president?")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(*result)
}
