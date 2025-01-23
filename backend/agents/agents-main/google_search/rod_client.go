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

func (r *RodClient) HtmlFromURL(url string, actions ...Action) (*string, error) {
	return r.htmlFromURL(url, actions...)
}

func (r *RodClient) htmlFromURL(url string, actions ...Action) (*string, error) {
	html := ""
	var error error
	err := rod.Try(func() {
		// Create page with URL directly, like in HtmlFromUrlCloseCookies
		page := r.browser.MustPage(url)
		defer page.Close()

		for _, action := range actions {
			switch action.Type() {
			case ActionTypeClick:
				clickAction := action.(*ClickAction)
				page.MustElement(clickAction.Element).MustClick()
			}
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
