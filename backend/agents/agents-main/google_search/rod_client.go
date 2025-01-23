package googleSearch

import (
	"fmt"
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

func NewRodClient(browser ...*rod.Browser) *RodClient {
	var b *rod.Browser
	if len(browser) > 0 {
		b = browser[0]
	} else {
		b = rod.New().MustConnect()
	}

	rc := &RodClient{
		browser: b,
	}

	err := rc.setup()
	if err != nil {
		panic(err)
	}

	return rc
}

func (r *RodClient) setup() error {
	_, err := r.htmlFromURL("https://www.google.com/maps", NewClickActionCloseGoogleCookies())
	return err
}

func NewNonHeadlessRodClient() *RodClient {
	l := launcher.New().Headless(false)

	return NewRodClient(rod.New().ControlURL(l.MustLaunch()).MustConnect())
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
			fmt.Printf("action: %+v\n", action)
			switch action.Type() {
			case ActionTypeClick:
				clickAction := action.(*ClickAction)

				if clickAction.Eq == nil {
					page.MustElement(clickAction.Element).MustClick()
					continue
				}

				elems := page.MustElements(clickAction.Element)
				if *clickAction.Eq >= len(elems) {
					error = fmt.Errorf("element not found, len == %d", len(elems))
					return
				}
				elems[*clickAction.Eq].MustClick()

			case ActionTypeWait:
				waitAction := action.(*WaitAction)
				if waitAction.Element != nil {
					page.MustElement(*waitAction.Element).MustWaitVisible()
				}
				if waitAction.Duration != nil {
					time.Sleep(*waitAction.Duration)
				}

			case ActionTypeNavigate:
				navigateAction := action.(*NavigateAction)
				page.MustNavigate(navigateAction.URL)
				page.MustWaitNavigation()
				page.MustWaitLoad()
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
