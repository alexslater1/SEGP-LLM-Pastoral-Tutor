package ragger

import "strings"

func cleanText(text string) string {
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
