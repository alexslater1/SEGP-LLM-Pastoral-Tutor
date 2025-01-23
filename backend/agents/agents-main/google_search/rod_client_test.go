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
