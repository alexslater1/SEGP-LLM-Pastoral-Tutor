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

func FormatUrl(uri string) (string, error) {
	// If no scheme is present, prepend "https://"
	if !regexp.MustCompile(`^[a-zA-Z]+://`).MatchString(uri) {
		uri = "https://" + uri
	}

	parsedUrl, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %s: %w", uri, err)
	}

	if parsedUrl.Host == "" {
		return "", fmt.Errorf("malformed url: %s", uri)
	}

	return parsedUrl.String(), nil
}
