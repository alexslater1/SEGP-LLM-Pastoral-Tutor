package googleSearch

type GoogleSearchClient interface {
	HtmlFromQuery(query string) (*string, error)
	HtmlFromURL(url string) (*string, error)
	HtmlFromUrlCloseCookies(url string) (*string, error)
}
