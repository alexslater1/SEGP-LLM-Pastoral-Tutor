package googleSearch

import "fmt"

type MockGoogleSearchClient struct {
	// Maps to store expected responses for different inputs
	queryResponses map[string]*string
	urlResponses   map[string]*string
}

// NewMockGoogleSearchClient creates a new MockGoogleSearchClient instance
func NewMockGoogleSearchClient() *MockGoogleSearchClient {
	return &MockGoogleSearchClient{
		queryResponses: make(map[string]*string),
		urlResponses:   make(map[string]*string),
	}
}

// SetQueryResponse sets the expected HTML response for a given search query
func (m *MockGoogleSearchClient) SetQueryResponse(query string, response string) {
	m.queryResponses[query] = &response
}

// SetURLResponse sets the expected HTML response for a given URL
func (m *MockGoogleSearchClient) SetURLResponse(url string, response string) {
	m.urlResponses[url] = &response
}

// HtmlFromQuery implements the GoogleSearchClient interface
func (m *MockGoogleSearchClient) HtmlFromQuery(query string, actions ...Action) (*string, error) {
	if response, ok := m.queryResponses[query]; ok {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response set for query: %s", query)
}

// HtmlFromURL implements the GoogleSearchClient interface
func (m *MockGoogleSearchClient) HtmlFromURL(url string, actions ...Action) (*string, error) {
	if response, ok := m.urlResponses[url]; ok {
		return response, nil
	}
	return nil, fmt.Errorf("no mock response set for URL: %s", url)
}
