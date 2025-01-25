package tools

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
	googleSearch "github.com/segp/agents-main/google_search"
)

type GoogleSearchResultsTool struct {
	googleSearchClient googleSearch.GoogleSearchClient
}

func NewGoogleSearchResultsTool(googleSearchClient googleSearch.GoogleSearchClient) *GoogleSearchResultsTool {

	return &GoogleSearchResultsTool{
		googleSearchClient: googleSearchClient,
	}
}

func (g *GoogleSearchResultsTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "google_search_results",
		Description: "Get the top google search results (urls and brief descriptions of the page content of that url) for a given query.",
		Parameters: []Parameter{
			{Name: "query", Description: "The query to search for top results for", Type: ParameterTypeString},
		},
	}
}

func (g *GoogleSearchResultsTool) GoogleSearchResultsFor(query string) (*string, error) {
	html, err := g.googleSearchClient.HtmlFromQuery(query, googleSearch.NewWaitElementAction("li.b_algo"))
	if err != nil {
		return nil, err
	}

	results, err := parseGoogleSearchResults(html)
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	jsonStr := string(json)
	return &jsonStr, nil
}

type GoogleSearchResult struct {
	URL         string
	Description string
}

func parseGoogleSearchResults(html *string) ([]GoogleSearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	results := []GoogleSearchResult{}
	doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
		description := s.Find("div.b_caption").Text()
		h2 := s.Find("h2")
		link, exists := h2.Find("a").Attr("href")
		if exists {
			results = append(results, GoogleSearchResult{
				URL:         link,
				Description: cleanString(description),
			})
		}
	})

	filteredResults := []GoogleSearchResult{}
	for _, result := range results {
		if !strings.HasSuffix(result.URL, ".pdf") && !strings.HasSuffix(result.URL, ".xml") {
			filteredResults = append(filteredResults, result)
		}
	}

	return filteredResults, nil
}
