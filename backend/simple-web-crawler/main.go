package main

import (
	"fmt"
)

// import (
// 	"fmt"
// 	"github.com/gocolly/colly/v2"
// 	"strings"
// )

// func main() {
// 	// Create a new collector
// 	c := colly.NewCollector(
// 		// Only allow imperial.ac.uk domains
// 		colly.AllowedDomains("www.imperial.ac.uk", "imperial.ac.uk"),
// 		// Set max depth to 3
// 		colly.MaxDepth(3),
// 	)

// 	// Create a map to store unique URLs
// 	visitedURLs := make(map[string]bool)

// 	// On every a element which has href attribute
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
// 		// Skip if current page URL doesn't contain /computing

// 		link := e.Attr("href")

// 		// Clean and normalize the URL
// 		link = strings.TrimSpace(link)
// 		if strings.HasPrefix(link, "/") {
// 			link = "https://www.imperial.ac.uk" + link
// 		}

// 		// Skip if not an imperial.ac.uk URL
// 		if !strings.Contains(link, "imperial.ac.uk") {
// 			return
// 		}

// 		// Skip if already visited
// 		if visitedURLs[link] {
// 			return
// 		}

// 		// Mark as visited and print
// 		visitedURLs[link] = true
// 		fmt.Printf("Found link: %s\n", link)

// 		// Visit the link
// 		e.Request.Visit(link)
// 	})

// 	// Before making a request
// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Printf("Visiting: %s\n", r.URL)
// 	})

// 	// Handle error
// 	c.OnError(func(r *colly.Response, err error) {
// 		fmt.Printf("Error on %s: %s\n", r.Request.URL, err)
// 	})

// 	// Start scraping
// 	c.Visit("https://www.imperial.ac.uk/placements/the-inplace-system/")

// 	// Print total number of unique URLs found
// 	fmt.Printf("\nTotal unique URLs found: %d\n", len(visitedURLs))

// 	fmt.Println("[")
// 	for url := range visitedURLs {
// 		fmt.Printf("\"%s\",\n", url)
// 	}
// 	fmt.Println("]")
// }

func main() {
	urls := filterDisformedUrls(filterDuplicates(urls))

	fmt.Println(len(urls))
	fmt.Println("[")
	for _, u := range urls {
		fmt.Printf("\"%s\",\n", u)
	}
	fmt.Println("]")
}
