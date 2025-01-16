package tools

import (
	"encoding/json"
	"strings"

	"html"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	googleSearch "github.com/segp/agents-main/google_search"
)

type GoogleSearchResultsTool struct {
	googleSearchClient googleSearch.GoogleSearchClient
}

type GoogleSearchResultsToolArgs struct {
	Query string
}

func NewGoogleSearchResultsTool(googleSearchClient googleSearch.GoogleSearchClient) Tool[GoogleSearchResultsToolArgs] {

	return &GoogleSearchResultsTool{
		googleSearchClient: googleSearchClient,
	}
}

func (g *GoogleSearchResultsTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "Google Search Results",
		Description: "Get the top google search results (urls and breif descriptions) for a given query.",
		Parameters: []Parameter{
			{Name: "query", Description: "The query to search for top results for", Type: ParameterTypeString},
		},
	}
}

func (g *GoogleSearchResultsTool) Call(args GoogleSearchResultsToolArgs) (*string, error) {
	html, err := g.googleSearchClient.HtmlFromQuery(args.Query)
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

func cleanString(s string) string {
	// Decode HTML entities
	s = html.UnescapeString(s)
	// Replace non-breaking spaces with regular spaces
	s = strings.ReplaceAll(s, "\u00a0", " ")
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)
	// Replace multiple spaces with a single space
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return s
}

func parseGoogleSearchResults(html *string) ([]GoogleSearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	results := []GoogleSearchResult{}
	doc.Find("div.MjjYud").Each(func(i int, s *goquery.Selection) {
		description := s.Find("div.VwiC3b.yXK7lf.p4wth.r025kc.hJNv6b").Text()
		link, exists := s.Find("a[jsname='UWckNb']").Attr("href")
		if exists {
			results = append(results, GoogleSearchResult{
				URL:         link,
				Description: cleanString(description),
			})
		}
	})

	return results, nil
}
