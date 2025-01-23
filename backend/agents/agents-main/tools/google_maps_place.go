package tools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	googleSearch "github.com/segp/agents-main/google_search"
)

type GoogleMapsPlaceTool struct {
	googleSearchClient googleSearch.GoogleSearchClient
}

func NewGoogleMapsPlaceTool(googleSearchClient googleSearch.GoogleSearchClient) *GoogleMapsPlaceTool {
	return &GoogleMapsPlaceTool{googleSearchClient: googleSearchClient}
}

func (g *GoogleMapsPlaceTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "google_maps_place",
		Description: "Returns lots of information about a specific place on Google Maps, given a GoogleMapsPlaceURL.",
		Parameters: []Parameter{
			{Name: "google_maps_place_url", Type: "string", Description: "The URL of the place on Google Maps. Comes from a previous google_maps_results tool call."},
		},
	}
}

func (g *GoogleMapsPlaceTool) PlaceDetailsFor(googleMapsPlaceURL string) (*string, error) {
	placeDetails, err := g.placeDetailsFor(googleMapsPlaceURL)
	if err != nil {
		return nil, err
	}

	jsonPlaceDetails, err := json.Marshal(placeDetails)
	if err != nil {
		return nil, err
	}

	jsonPlaceDetailsString := string(jsonPlaceDetails)

	return &jsonPlaceDetailsString, nil
}

func (g *GoogleMapsPlaceTool) placeDetailsFor(googleMapsPlaceURL string) (*PlaceDetails, error) {
	placeOverview, err := g.placeOverviewFor(googleMapsPlaceURL)
	if err != nil {
		return nil, err
	}

	return &PlaceDetails{Overview: *placeOverview}, nil
}

func (g *GoogleMapsPlaceTool) placeOverviewFor(googleMapsPlaceURL string) (*PlaceOverview, error) {
	if err := verifyIsGoogleMapsPlaceURL(googleMapsPlaceURL); err != nil {
		return nil, err
	}

	actions := []googleSearch.Action{
		googleSearch.NewClickActionCloseGoogleCookies(),
		googleSearch.NewWaitElementAction("h1.DUwDvf.lfPIob"),
	}

	html, err := g.googleSearchClient.HtmlFromURL(googleMapsPlaceURL, actions...)
	if err != nil {
		return nil, err
	}

	return parsePlaceOverview(html)
}

func (g *GoogleMapsPlaceTool) placeAboutFor(googleMapsPlaceURL string) (*PlaceAbout, error) {
	actions := []googleSearch.Action{
		googleSearch.NewClickActionCloseGoogleCookies(),
		googleSearch.NewClickAction("button.hh2c6.G7m0Af").Nth(2),
		googleSearch.NewWaitElementAction("div.m6QErb.DxyBCb.kA9KIf.dS8AEf.XiKgde.d2JYHf"),
	}

	html, err := g.googleSearchClient.HtmlFromURL(googleMapsPlaceURL, actions...)
	if err != nil {
		return nil, err
	}

	return parsePlaceAbouts(html)
}

type PlaceDetails struct {
	Overview PlaceOverview
	About    PlaceAbout
}

type PlaceAbout struct {
	Description string
	Items       []PlaceAboutSection
}

type PlaceAboutSection struct {
	Title string
	Items []PlaceAboutSectionItem
}

type PlaceAboutSectionItem struct {
	Item        string
	IsAvailable bool
}

type PlaceOverview struct {
	Name           string
	Stars          string
	Type           string
	Description    string
	Details        string
	ReviewPreviews []string
	Reviews        []Review
}

type Review struct {
	Description string
	TimeAgo     string
}

func parsePlaceAbouts(html *string) (*PlaceAbout, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	placeAboutSections := []PlaceAboutSection{}
	doc.Find("div.iP2t7d.fontBodyMedium").Each(func(i int, s *goquery.Selection) {
		placeAboutSectionItems := []PlaceAboutSectionItem{}
		s.Find("li.hpLkke").Each(func(i int, s *goquery.Selection) {
			placeAboutSectionItems = append(placeAboutSectionItems, PlaceAboutSectionItem{
				Item:        cleanString(s.Find("span").Eq(1).Text()),
				IsAvailable: s.Find("f5BGzb.google-symbols").HasClass("SwaGS"),
			})
		})

		placeAboutSections = append(placeAboutSections, PlaceAboutSection{
			Title: cleanString(s.Find("h2").Text()),
			Items: placeAboutSectionItems,
		})
	})

	return &PlaceAbout{
		Description: cleanString(doc.Find("div.PbZDve p.fontBodyMedium span.HlvSq").Text()),
		Items:       placeAboutSections,
	}, nil
}

func parsePlaceOverview(html *string) (*PlaceOverview, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		return nil, err
	}

	reviewPreviews := []string{}
	doc.Find("div.tBizfc").Each(func(i int, s *goquery.Selection) {
		reviewPreviews = append(reviewPreviews, s.Text())
	})

	reviews := []Review{}
	doc.Find("div.jftiEf").Each(func(i int, s *goquery.Selection) {
		reviews = append(reviews, Review{
			Description: cleanString(s.Find("span.wiI7pd").Text()),
			TimeAgo:     cleanString(s.Find("span.rsqaWe").Text()),
		})
	})

	detailsStr := ""
	doc.Find("div.Io6YTe.fontBodyMedium.kR99db.fdkmkc ").Each(func(i int, s *goquery.Selection) {
		detailsStr += cleanString(s.Text()) + " \n"
	})

	doc.Find("div.RcCsl.w4vB1d.NOE9ve.M0S7ae.AG25L").Each(func(i int, s *goquery.Selection) {
		detailsStr += cleanString(s.Find("div.Io6YTe.fontBodyMedium.kR99db.fdkmkc ").Text()) + ":" + cleanString(s.Find("a.CsEnBe").AttrOr("href", "")) + " \n"
	})

	placeOverview := &PlaceOverview{
		Name:           cleanString(doc.Find("h1.DUwDvf.lfPIob").Text()),
		Stars:          cleanString(doc.Find("div.F7nice span span").First().Text()),
		Type:           cleanString(doc.Find("button.DkEaL").Text()),
		Description:    cleanString(doc.Find("div.PYvSYb").Text()),
		Details:        detailsStr,
		ReviewPreviews: reviewPreviews,
	}

	return placeOverview, nil
}

func verifyIsGoogleMapsPlaceURL(URL string) error {
	if URL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check if it's a Google Maps URL
	if parsedURL.Host != "google.com" && parsedURL.Host != "www.google.com" {
		return fmt.Errorf("URL must be from google.com")
	}

	// Check if it's a place URL
	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 3 || pathParts[1] != "maps" || pathParts[2] != "place" {
		return fmt.Errorf("URL must be a Google Maps place URL (format: https://google.com/maps/place/...)")
	}

	return nil
}
