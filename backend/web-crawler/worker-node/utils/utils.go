package utils

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func Required[T any](value T, name string) T {
	if reflect.ValueOf(value).IsZero() {
		panic(fmt.Sprintf("%s is required", name))
	}
	return value
}

func RequiredInt(value string, name string) int {
	if value == "" {
		log.Fatalf("%s environment variable is required", name)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Failed to parse %s as integer: %v", name, err)
	}
	return intValue
}

func HtmlToMarkdown(html *string) (string, error) {
	// Create a new converter with default options
	converter := md.NewConverter("", true, nil)

	// Convert HTML to Markdown
	markdown, err := converter.ConvertString(*html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func CleanText(text string) string {
	re := regexp.MustCompile(`\n\s*\n`)
	text = re.ReplaceAllString(text, "\n\n")
	return text
}

func IsValidUrl(uri string) bool {
	parsedUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	// Check if the URL has a host
	return parsedUrl.Host != ""
}
