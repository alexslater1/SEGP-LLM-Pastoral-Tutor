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

func TestRodClientClickNth(t *testing.T) {
	actions := []Action{
		NewClickActionCloseGoogleCookies(),
		NewWaitElementAction("div.yx21af.lLU2pe.XDi3Bc"),
		NewClickAction("button.hh2c6").Nth(2),
	}

	url := "https://www.google.com/maps/place/Boucherie+Union+Square/@40.7372552,-73.9882246,17z/data=!3m1!5s0x89c259a1e966311b:0xaa9b5a583f466ed2!4m15!1m8!3m7!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!2sBoucherie+Union+Square!8m2!3d40.7372552!4d-73.9882246!10e1!16s%2Fg%2F11hbv5rh0_!3m5!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!8m2!3d40.7372552!4d-73.9882246!16s%2Fg%2F11hbv5rh0_?authuser=0&hl=en&entry=ttu&g_ep=EgoyMDI1MDEyMC4wIKXMDSoASAFQAw%3D%3D"

	client := NewNonHeadlessRodClient()
	html, err := client.HtmlFromURL(url, actions...)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*html)
}
