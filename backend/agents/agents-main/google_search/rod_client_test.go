package googleSearch

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	markdown "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
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

	url := "https://www.google.com/maps/place/Boucherie+Union+Square/data=!4m7!3m6!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!8m2!3d40.7372552!4d-73.9882246!16s%2Fg%2F11hbv5rh0_!19sChIJc1Vf7KFZwokR1YL2Rn9oxi8?authuser=0&hl=en&rclk=1"
	actions := []Action{
		NewClickActionCloseGoogleCookies(),
		NewWaitElementAction("h1.DUwDvf.lfPIob"),
	}
	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url, actions...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}

func TestRodClientClickNth(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	actions := []Action{
		NewClickActionCloseGoogleCookies(),
		NewWaitElementAction("div.yx21af.lLU2pe.XDi3Bc"),
		NewClickAction("button.hh2c6").Nth(2),
	}

	url := "https://www.google.com/maps/place/Boucherie+Union+Square/data=!4m7!3m6!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!8m2!3d40.7372552!4d-73.9882246!16s%2Fg%2F11hbv5rh0_!19sChIJc1Vf7KFZwokR1YL2Rn9oxi8?authuser=0&hl=en&rclk=1"

	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url, actions...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}

func TestRodClientNavigate(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	url := "https://www.ethanhosier.com/"

	actions := []Action{
		NewWaitDurationAction(time.Second * 5),
		NewNavigateAction("https://www.wikipedia.org/"),
	}
	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url, actions...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}

func TestRodClientFromUrl2(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}
	
	url := "https://www.exchangerates.org.uk/US-Dollar-USD-currency-table.html"

	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("###html:")
	t.Log(*html)

	mdContent, err := relevantPageMarkdownContent(html)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("###mdContent:")
	t.Log(*mdContent)
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
