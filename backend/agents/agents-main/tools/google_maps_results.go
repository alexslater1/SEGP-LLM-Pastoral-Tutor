package tools

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	googleSearch "github.com/segp/agents-main/google_search"
)

type GoogleMapsResult struct {
	Title              string
	Rating             string
	GoogleMapsPlaceURL string
	Info               string
	Tags               string
}

type GoolgeMapsResultsTool struct {
	googleSearchClient googleSearch.GoogleSearchClient
}

func NewGoogleMapsResultsTool(gc googleSearch.GoogleSearchClient) *GoolgeMapsResultsTool {
	return &GoolgeMapsResultsTool{
		googleSearchClient: gc,
	}
}

func (g *GoolgeMapsResultsTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "google_maps_results",
		Description: "Returns the results of a google maps search. For each place listed, will get the title, rating, PLACES URL and then some info + tags about the place.",
		Parameters: []Parameter{
			{Name: "query", Description: "The query to search for results for", Type: ParameterTypeString},
		},
	}
}

func (g *GoolgeMapsResultsTool) GoogleMapsResultsFor(query string) (*string, error) {
	actions := []googleSearch.Action{
		googleSearch.NewClickAction(".UywwFc-LgbsSe.UywwFc-LgbsSe-OWXEXe-dgl2Hf.XWZjwc"),
	}
	html, err := g.googleSearchClient.HtmlFromURL("https://www.google.com/maps/search/"+url.QueryEscape(query), actions...)
	if err != nil {
		return nil, err
	}

	results, err := g.parseGoogleMapsResults(html)
	if err != nil {
		return nil, err
	}

	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	resultsStr := string(resultsBytes)

	return &resultsStr, nil
}

func (g *GoolgeMapsResultsTool) parseGoogleMapsResults(html *string) ([]GoogleMapsResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	results := []GoogleMapsResult{}
	doc.Find("div.Nv2PK.THOPZb.CpccDe").Each(func(i int, s *goquery.Selection) {
		results = append(results, GoogleMapsResult{
			Title:              cleanString(s.Find("div.NrDZNb").Text()),
			Rating:             cleanString(s.Find("span.MW4etd").Text()),
			Info:               cleanString(s.Find("div.W4Efsd").Text()),
			Tags:               cleanString(s.Find("div.n8sPKe.fontBodySmall.ccePVe").Text()),
			GoogleMapsPlaceURL: s.Find("a.hfpxzc").AttrOr("href", ""),
		})
	})

	return results, nil
}
