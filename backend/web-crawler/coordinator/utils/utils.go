package utils

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

func CleanText(text string) string {
	// Split into lines, trim each line, and handle multiple newlines
	lines := strings.Split(text, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}

	// Join with double newlines and clean up any remaining multiple newlines
	text = strings.Join(cleanedLines, "\n")
	re := regexp.MustCompile(`\n\s*\n`)
	text = re.ReplaceAllString(text, "\n\n")

	// Remove zero-width characters
	text = strings.ReplaceAll(text, "\u200c", "") // Remove zero-width non-joiner
	text = strings.ReplaceAll(text, "\u200b", "") // Remove zero-width space

	// Handle escape sequences
	text = strings.ReplaceAll(text, "\\n", "\n")
	text = strings.ReplaceAll(text, "\\\"", "\"")
	text = strings.ReplaceAll(text, "\\\\", "\\")

	// Ensure text doesn't end with a partial escape sequence
	text = strings.TrimSuffix(text, "\\")

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
