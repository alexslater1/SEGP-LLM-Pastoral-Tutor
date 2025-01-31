package utils

import (
	"errors"
	"testing"
)

func TestIsValidUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		err      error
	}{
		{
			name:     "valid http url",
			url:      "http://example.com",
			expected: "http://example.com",
			err:      nil,
		},
		{
			name:     "valid https url",
			url:      "https://example.com/path?query=value",
			expected: "https://example.com/path?query=value",
			err:      nil,
		},
		{
			name:     "valid url but missing scheme",
			url:      "example.com",
			expected: "https://example.com",
			err:      nil,
		},
		{
			name:     "invalid url - empty string",
			url:      "",
			expected: "",
			err:      errors.New("empty string"),
		},
		{
			name:     "invalid url - malformed",
			url:      "http://",
			expected: "",
			err:      errors.New("malformed url"),
		},
		{
			name:     "valid url - missing scheme",
			url:      "accounting-by-post.com",
			expected: "https://accounting-by-post.com",
			err:      nil,
		},
		{
			name:     "valid url - missing scheme",
			url:      "fdu.org.ua",
			expected: "https://fdu.org.ua",
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatUrl(tt.url)
			if err != nil && tt.err == nil {
				t.Errorf("FormatUrl(%q) Error = %v, want %v", tt.url, err, tt.expected)
			}

			if result != tt.expected {
				t.Errorf("FormatUrl(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}
