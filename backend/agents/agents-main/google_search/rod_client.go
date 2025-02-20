package googleSearch

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const (
	timeoutDuration = 15 * time.Second
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

func (r *RodClient) HtmlFromQuery(query string, actions ...Action) (*string, error) {
	encodedQuery := url.QueryEscape(query)
	url := "https://www.bing.com/search?form=&q=" + encodedQuery + "&form=QBLH&sp=-1&lq=0&pq=" + encodedQuery + "&sc=12-18&qs=n&sk=&cvid=CD45C64E103445449518EF0A5C1C3315&ghsh=0&ghacc=0&ghpl="

	fmt.Println("url", url)
	return r.htmlFromURL(url, actions...)
}

func (r *RodClient) HtmlFromURL(url string, actions ...Action) (*string, error) {
	return r.htmlFromURL(url, actions...)
}

func (r *RodClient) htmlFromURL(url string, actions ...Action) (*string, error) {
	type result struct {
		html string
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		html := ""
		var customErr error
		err := rod.Try(func() {
			page := r.browser.MustPage(url)
			defer page.Close()

			// Add initial wait for page load
			page.MustWaitLoad()

			for _, action := range actions {
				switch action.Type() {
				case ActionTypeClick:
					clickAction := action.(*ClickAction)

					if clickAction.Eq == nil {
						page.MustElement(clickAction.Element).MustClick()
						continue
					}

					elems := page.MustElements(clickAction.Element)
					if *clickAction.Eq >= len(elems) {
						customErr = fmt.Errorf("element not found, len == %d", len(elems))
						return
					}
					elems[*clickAction.Eq].MustClick()

				case ActionTypeWait:
					waitAction := action.(*WaitAction)
					if waitAction.Element != nil {
						err := rod.Try(func() {
							if waitAction.Timeout == nil {
								page.MustElement(*waitAction.Element).MustWaitVisible()
							} else {
								page.MustElement(*waitAction.Element).Timeout(*waitAction.Timeout).MustWaitVisible()
							}
						})
						if err != nil {
							log.Println("Warning: timeout waiting for element")
							continue
						}
					}
					if waitAction.Duration != nil {
						time.Sleep(*waitAction.Duration)
					}

				case ActionTypeNavigate:
					navigateAction := action.(*NavigateAction)
					page.MustNavigate(navigateAction.URL)
					page.MustWaitNavigation()
				}
				page.MustWaitLoad()
			}

			// Increase timeout for request idle
			page.WaitRequestIdle(5*time.Second, []string{""}, []string{}, excludeTypes)

			h, err := page.HTML()
			if err != nil {
				customErr = err
				return
			}
			html = h
		})

		// If rod.Try returned an error, return a safe error HTML.
		if err != nil {
			log.Println("Warning: there was an error in rod client. Returning error html but continuing.")
			resultChan <- result{
				html: fmt.Sprintf(`<html><body><h1>Error accessing page %s</h1></body></html>`, url),
				err:  nil,
			}
			return
		}
		if customErr != nil {
			resultChan <- result{
				html: fmt.Sprintf(`<html><body><h1>Error accessing page %s: %s</h1></body></html>`, url, customErr.Error()),
				err:  nil,
			}
			return
		}
		resultChan <- result{
			html: html,
			err:  nil,
		}
	}()

	// Wait for the result or timeout after 15 seconds.
	select {
	case res := <-resultChan:
		return &res.html, res.err
	case <-time.After(timeoutDuration):
		timeoutHTML := `<html><body><h1>Timeout accessing page</h1></body></html>`
		return &timeoutHTML, nil
	}
}
