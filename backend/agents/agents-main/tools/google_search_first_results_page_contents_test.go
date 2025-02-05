package tools

import (
	"fmt"
	"os"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
)

func TestGoogleSearchFirstResultsPageContents(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}
	
	tool := NewGoogleSearchFirstResultsPageContentsTool(googleSearch.NewNonHeadlessRodClient(), 3)

	results, err := tool.GoogleSearchFirstResultsPageContentsFor("usd price")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(*results)
}
