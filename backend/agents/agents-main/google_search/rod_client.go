package googleSearch

import (
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

var (
	// Exclude non-text resources like images, videos, etc.
	excludeTypes = []proto.NetworkResourceType{
		proto.NetworkResourceTypeImage,
		proto.NetworkResourceTypeMedia,
		proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeStylesheet,
	}
)

// TODO: add page pool
type RodClient struct {
	browser *rod.Browser
}

func NewRodClient() *RodClient {

	return &RodClient{
		browser: rod.New().MustConnect(),
	}
}

func NewNonHeadlessRodClient() *RodClient {
	l := launcher.New().Headless(false)

	return &RodClient{
		browser: rod.New().ControlURL(l.MustLaunch()).MustConnect(),
	}
}

func (r *RodClient) HtmlFromQuery(query string) (*string, error) {
	encodedQuery := url.QueryEscape(query)
	url := "https://www.bing.com/search?form=&q=" + encodedQuery
	return r.htmlFromURL(url)
}

func (r *RodClient) HtmlFromURL(url string) (*string, error) {
	return r.htmlFromURL(url)
}

func (r *RodClient) htmlFromURL(url string) (*string, error) {
	html := ""
	var error error
	err := rod.Try(func() {
		page := r.browser.MustPage()
		defer page.Close()

		// Configure page to look like a regular browser
		page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
			Platform:  "Windows",
		})

		err := page.Navigate(url)
		if err != nil {
			error = err
			return
		}

		// Wait for the page to fully load
		page.MustWaitLoad()

		page.WaitRequestIdle(time.Second*3, []string{""}, []string{}, excludeTypes)

		h, err := page.HTML()
		if err != nil {
			error = err
			return
		}

		html = h
	})
	if err != nil {
		return nil, err
	}

	return &html, error
}

func (r *RodClient) HtmlFromUrlCloseCookies(url string) (*string, error) {
	html := ""
	var error error
	err := rod.Try(func() {
		page := r.browser.MustPage(url)
		defer page.Close()

		button := page.MustElement(".UywwFc-LgbsSe.UywwFc-LgbsSe-OWXEXe-dgl2Hf.XWZjwc")
		button.MustClick()

		// Wait for the page to fully load
		page.MustWaitLoad()

		page.WaitRequestIdle(time.Second*3, []string{""}, []string{}, excludeTypes)

		h, err := page.HTML()
		if err != nil {
			error = err
			return
		}

		html = h
	})
	if err != nil {
		return nil, err
	}

	return &html, error
}
