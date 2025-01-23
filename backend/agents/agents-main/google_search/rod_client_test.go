package googleSearch

import (
	"os"
	"testing"
)

func TestRodClientHtmlFromQuery(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	query := "hello"
	client := NewRodClient()
	html, err := client.HtmlFromQuery(query)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}

func TestRodClientHtmlFromURL(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	url := "https://www.google.com/maps"
	actions := []Action{
		NewClickAction(".UywwFc-LgbsSe.UywwFc-LgbsSe-OWXEXe-dgl2Hf.XWZjwc"),
	}
	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url, actions...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}
