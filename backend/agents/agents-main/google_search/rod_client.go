package googleSearch

import (
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

var (
	// Exclude non-text resources like images, videos, etc.
	excludeTypes = []proto.NetworkResourceType{
		proto.NetworkResourceTypeImage,
		proto.NetworkResourceTypeMedia,
		proto.NetworkResourceTypeFont,
		proto.NetworkResourceTypeStylesheet,
		proto.NetworkResourceTypeScript,
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

func (r *RodClient) HtmlFromQuery(query string) (*string, error) {
	encodedQuery := url.QueryEscape(query)
	url := "https://www.google.com/search?q=" + encodedQuery
	return r.htmlFromURL(url)
}

func (r *RodClient) HtmlFromURL(url string) (*string, error) {
	return r.htmlFromURL(url)
}

func (r *RodClient) htmlFromURL(url string) (*string, error) {
	html := ""
	var error error
	rod.Try(func() {
		page := r.browser.MustPage()
		defer page.Close()

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
	return &html, error
}
