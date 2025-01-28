package utils

import "testing"

func TestIsValidUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "valid http url",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "valid https url",
			url:      "https://example.com/path?query=value",
			expected: true,
		},
		{
			name:     "invalid url - missing scheme",
			url:      "example.com",
			expected: false,
		},
		{
			name:     "invalid url - empty string",
			url:      "",
			expected: false,
		},
		{
			name:     "invalid url - malformed",
			url:      "http://",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidUrl(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidUrl(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}
