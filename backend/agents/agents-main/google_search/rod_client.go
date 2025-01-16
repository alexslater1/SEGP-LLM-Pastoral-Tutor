package googleSearch

import (
	"net/url"
	"time"

	"github.com/go-rod/rod"
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
	page := r.browser.MustPage()
	defer page.Close()

	err := page.Navigate(url)
	if err != nil {
		return nil, err
	}

	err = page.WaitLoad()
	if err != nil {
		return nil, err
	}

	err = page.WaitStable(time.Second * 2)
	if err != nil {
		return nil, err
	}

	html, err := page.HTML()
	if err != nil {
		return nil, err
	}
	return &html, nil
}
