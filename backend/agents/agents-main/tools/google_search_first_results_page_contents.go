package tools

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/utils"
)

type GoogleSearchFirstResultsPageContentsTool struct {
	googleSearchClient googleSearch.GoogleSearchClient

	numberOfResults int
}

func NewGoogleSearchFirstResultsPageContentsTool(googleSearchClient googleSearch.GoogleSearchClient, numberOfResults int) *GoogleSearchFirstResultsPageContentsTool {
	return &GoogleSearchFirstResultsPageContentsTool{
		googleSearchClient: googleSearchClient,
		numberOfResults:    numberOfResults,
	}
}

func (g *GoogleSearchFirstResultsPageContentsTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "google_search_first_results_page_contents",
		Description: fmt.Sprintf("Returns the contents of the first %d search results for a query.", g.numberOfResults),
		Parameters: []Parameter{
			{Name: "query", Type: "string", Description: "The query to search for."},
		},
	}
}

func (g *GoogleSearchFirstResultsPageContentsTool) GoogleSearchFirstResultsPageContentsFor(query string) (*string, error) {
	results, err := g.googleSearchClient.HtmlFromQuery(query, googleSearch.NewWaitElementAction("li.b_algo").WithTimeout(6*time.Second))
	if err != nil {
		return nil, err
	}

	parsedResults, err := parseGoogleSearchResults(results)
	if err != nil {
		return nil, err
	}

	topNResults := parsedResults[:int(math.Min(float64(g.numberOfResults), float64(len(parsedResults))))]
	topNUrls := []string{}
	for _, result := range topNResults {
		topNUrls = append(topNUrls, result.URL)
	}

	htmlTasks := utils.DoAsyncList(topNUrls, func(url string) (*string, error) {
		time.Sleep(time.Duration(rand.Float64()) * time.Second)
		return g.googleSearchClient.HtmlFromURL(url)
	})

	pageContents := make([]string, len(htmlTasks))

	for i, ht := range htmlTasks {
		html, err := ht.Get()
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
		if err != nil {
			return nil, err
		}

		// Remove unwanted elements
		doc.Find("script, style, noscript, iframe, header, footer, nav, aside").Remove()

		// Try to find main content area in order of likelihood
		var content string
		selectors := []string{
			"main", "article",
			"[role='main']",
			".main-content", ".content-main", "#content",
			".post-content", ".article-content",
		}

		for _, selector := range selectors {
			selection := doc.Find(selector)
			if selection.Length() > 0 {
				content = utils.CleanText(selection.First().Text())
				break
			}
		}

		// Fallback to body if no main content area found
		if content == "" {
			content = utils.CleanText(doc.Find("body").First().Text())
		}

		// Limit content length (e.g., first 1000 characters)
		const maxChars = 1000
		if len(content) > maxChars {
			content = content[:maxChars] + "..."
		}

		pageContents[i] = content + "\n"
	}

	pageContentsStr := strings.Join(pageContents, "\n\n")
	return &pageContentsStr, nil
}
