package scraper

import (
	"testing"
)

func TestMockScraper(t *testing.T) {
	t.Run("HtmlFrom returns set content", func(t *testing.T) {
		mock := NewMockScraper()
		expectedContent := "<html><body>Test Content</body></html>"
		testURL := "http://example.com"

		mock.SetHtmlContent(testURL, expectedContent)

		result, err := mock.HtmlFrom(testURL)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if *result != expectedContent {
			t.Errorf("Expected content %q, got %q", expectedContent, *result)
		}
	})

	t.Run("HtmlFrom returns error for unset URL", func(t *testing.T) {
		mock := NewMockScraper()
		testURL := "http://example.com"

		_, err := mock.HtmlFrom(testURL)
		if err == nil {
			t.Error("Expected error for unset URL, got nil")
		}
	})

	t.Run("HtmlFromTag returns set content", func(t *testing.T) {
		mock := NewMockScraper()
		expectedContent := "<div>Test Content</div>"
		testURL := "http://example.com"

		mock.SetHtmlContent(testURL, expectedContent)

		result, err := mock.HtmlFromTag(testURL, "div")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if *result != expectedContent {
			t.Errorf("Expected content %q, got %q", expectedContent, *result)
		}
	})

	t.Run("HtmlFromTag returns error for unset URL", func(t *testing.T) {
		mock := NewMockScraper()
		testURL := "http://example.com"

		_, err := mock.HtmlFromTag(testURL, "div")
		if err == nil {
			t.Error("Expected error for unset URL, got nil")
		}
	})
}
