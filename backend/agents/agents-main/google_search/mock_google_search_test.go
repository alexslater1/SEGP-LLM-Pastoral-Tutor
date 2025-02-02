package googleSearch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMockGoogleSearchClient(t *testing.T) {
	t.Run("HtmlFromQuery", func(t *testing.T) {
		mock := NewMockGoogleSearchClient()
		expectedHTML := "<html><body>Test Result</body></html>"
		mock.SetQueryResponse("test query", expectedHTML)

		// Test successful case
		response, err := mock.HtmlFromQuery("test query")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedHTML, *response)

		// Test with actions (should not affect the result)
		response, err = mock.HtmlFromQuery("test query",
			NewClickAction("button"),
			NewWaitAction("element", time.Second),
		)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedHTML, *response)

		// Test error case for unknown query
		response, err = mock.HtmlFromQuery("unknown query")
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "no mock response set for query")
	})

	t.Run("HtmlFromURL", func(t *testing.T) {
		mock := NewMockGoogleSearchClient()
		expectedHTML := "<html><body>Page Content</body></html>"
		mock.SetURLResponse("https://example.com", expectedHTML)

		// Test successful case
		response, err := mock.HtmlFromURL("https://example.com")
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedHTML, *response)

		// Test with actions (should not affect the result)
		response, err = mock.HtmlFromURL("https://example.com",
			NewNavigateAction("https://example.com/page2"),
			NewWaitElementAction("element"),
		)
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, expectedHTML, *response)

		// Test error case for unknown URL
		response, err = mock.HtmlFromURL("https://unknown.com")
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "no mock response set for URL")
	})
}
