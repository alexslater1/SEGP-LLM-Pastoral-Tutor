package scraper

type Scraper interface {
	HtmlFrom(url string) (*string, error)
	HtmlFromTag(url string, tag string) (*string, error)
}
