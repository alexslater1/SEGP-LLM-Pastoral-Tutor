package tools

import googleSearch "github.com/segp/agents-main/google_search"
import (
	"bytes"
	"fmt"
	markdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type GoogleSearchUrlTool struct {
	googleSearchClient googleSearch.GoogleSearchClient
}

type GoogleSearchUrlToolArgs struct {
	URL string
}

func NewGoogleSearchUrlTool(googleSearchClient googleSearch.GoogleSearchClient) Tool[GoogleSearchUrlToolArgs] {
	return &GoogleSearchUrlTool{
		googleSearchClient: googleSearchClient,
	}
}

func (g *GoogleSearchUrlTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "Google Search URL",
		Description: "Search a specific url, and return the scraped content of that page",
		Parameters: []Parameter{
			{Name: "URL", Type: "string", Description: "The URL to search."},
		},
	}
}

func (g *GoogleSearchUrlTool) Call(args GoogleSearchUrlToolArgs) (*string, error) {
	if err := validateURL(args.URL); err != nil {
		return nil, err
	}

	html, err := g.googleSearchClient.HtmlFromURL(args.URL)
	if err != nil {
		return nil, err
	}

	return relevantPageMarkdownContent(html)
}

func validateURL(urlStr string) error {
	if !strings.HasPrefix(urlStr, "http") {
		return fmt.Errorf("url must start with http or https")
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("url must start with http or https")
	}

	return nil
}

func relevantPageMarkdownContent(html *string) (*string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Extract the relevant content
	var buffer bytes.Buffer
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		htmlContent, err := s.Html()
		if err == nil {
			buffer.WriteString(htmlContent)
		}
	})

	// Convert the extracted HTML content to markdown
	converter := markdown.NewConverter("", true, nil)
	markdownStr, err := converter.ConvertString(buffer.String())
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	return &markdownStr, nil
}
