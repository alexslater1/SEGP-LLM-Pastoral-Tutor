package tools

import (
	"html"
	"regexp"
	"strings"
)

func cleanString(s string) string {
	// Decode HTML entities
	s = html.UnescapeString(s)
	// Replace non-breaking spaces with regular spaces
	s = strings.ReplaceAll(s, "\u00a0", " ")
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)
	// Replace multiple spaces with a single space
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	return s
}
