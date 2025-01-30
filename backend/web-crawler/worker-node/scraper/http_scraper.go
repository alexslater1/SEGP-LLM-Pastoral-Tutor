package scraper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethanhosier/worker-node/utils"
)

type HttpScraper struct {
}

func NewHttpScraper() *HttpScraper {
	return &HttpScraper{}
}

func (h *HttpScraper) HtmlFrom(url string) (*string, error) {

	formattedUrl, err := utils.FormatUrl(url)
	if err != nil {
		return nil, fmt.Errorf("failed to format url %s: %w", url, err)
	}

	// Make HTTP GET request
	resp, err := http.Get(formattedUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to make http get request for url %s: %w", formattedUrl, err)
	}
	defer resp.Body.Close()

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Get the HTML content
	html, err := doc.Find("html").Html()
	if err != nil {
		return nil, err
	}
	return &html, nil
}

func (h *HttpScraper) HtmlFromTag(url string, tag string) (*string, error) {
	formattedUrl, err := utils.FormatUrl(url)
	if err != nil {
		return nil, fmt.Errorf("failed to format url %s: %w", url, err)
	}

	// Make HTTP GET request directly (don't reuse HtmlFrom to avoid double parsing)
	resp, err := http.Get(formattedUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a goquery document from the HTTP response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Get the HTML content for the specific tag
	html, err := doc.Find(tag).Html()
	if err != nil {
		return nil, err
	}
	return &html, nil
}
