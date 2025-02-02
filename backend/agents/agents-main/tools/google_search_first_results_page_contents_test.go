package tools

import (
	"fmt"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
)

func TestGoogleSearchFirstResultsPageContents(t *testing.T) {
	tool := NewGoogleSearchFirstResultsPageContentsTool(googleSearch.NewNonHeadlessRodClient(), 3)

	results, err := tool.GoogleSearchFirstResultsPageContentsFor("usd price")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(*results)
}
