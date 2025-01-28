package scraper

import "fmt"

type MockScraper struct {
	htmlContent map[string]string
}

func NewMockScraper() *MockScraper {
	return &MockScraper{
		htmlContent: make(map[string]string),
	}
}

// SetHtmlContent sets the HTML content to be returned for a specific URL
func (m *MockScraper) SetHtmlContent(url string, content string) {
	m.htmlContent[url] = content
}

func (m *MockScraper) HtmlFrom(url string) (*string, error) {
	if content, exists := m.htmlContent[url]; exists {
		return &content, nil
	}
	return nil, fmt.Errorf("no mock content set for URL: %s", url)
}

func (m *MockScraper) HtmlFromTag(url string, tag string) (*string, error) {
	// For simplicity, we'll return the entire mock content
	// In a more sophisticated implementation, you could store and return tag-specific content
	if content, exists := m.htmlContent[url]; exists {
		return &content, nil
	}
	return nil, fmt.Errorf("no mock content set for URL: %s", url)
}
