package scraper

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestHttpScraper(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping HTTP scraper in CI/CD")
	}

	scraper := NewHttpScraper()

	html, err := scraper.HtmlFrom("https://www.imperial.ac.uk/study/courses/undergraduate/computing-meng/")
	if err != nil {
		t.Fatalf("Error getting HTML: %v", err)
	}

	t.Logf("HTML: %v", *html)
}

func TestHttpScraperHtmlFromTag(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping HTML from tag in CI/CD")
	}

	scraper := NewHttpScraper()

	html, err := scraper.HtmlFromTag("https://www.imperial.ac.uk/study/courses/undergraduate/computing-meng/", "main")
	if err != nil {
		t.Fatalf("Error getting HTML: %v", err)
	}

	t.Logf("HTML: %v", *html)
}
