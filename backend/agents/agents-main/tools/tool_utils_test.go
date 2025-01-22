package tools

import "testing"

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML entities",
			input:    "Hello &amp; World",
			expected: "Hello & World",
		},
		{
			name:     "Non-breaking spaces",
			input:    "Hello\u00a0World",
			expected: "Hello World",
		},
		{
			name:     "Multiple spaces",
			input:    "Hello    World   !",
			expected: "Hello World !",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  Hello World  ",
			expected: "Hello World",
		},
		{
			name:     "Combined cases",
			input:    "  Hello\u00a0&amp;\u00a0World  ",
			expected: "Hello & World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanString(tt.input)
			if got != tt.expected {
				t.Errorf("cleanString() = %q, want %q", got, tt.expected)
			}
		})
	}
}
