package tools

import (
	"fmt"
	"os"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
)

func TestParseGoogleMapsResults(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in non-CI/CD environment")
	}

	tool := NewGoogleMapsResultsTool(googleSearch.NewNonHeadlessRodClient())

	results, err := tool.GoogleMapsResultsFor("restaurant in new york")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(*results)
}
