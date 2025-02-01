package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCleanText(t *testing.T) {
	text := `In this section

                                Study

                    Teacher reference guidance

                Study

                        Imperial Home
                        StudyApplyUndergraduateApplication processApplication reference

                Application reference

Application process

Teacher reference guidance

		Everything you need to know about your application reference

UCAS advice: How to get an undergraduate reference

You only need one reference for your UCAS application. It’s usually written for you by a teacher or tutor who knows you – if you’re applying through a school, college or centre registered with UCAS, they’ll do this before they send your completed application to UCAS; you won’t need to do anything for it. See below for more information about choosing a reference if you’re applying independently. What the reference is for
We really value the information provided by your referee as it helps us to gain additional context around your application, in particular about:

your school
your learning environment
the grades you are predicted to achieve

For example, the reference is an opportunity for your referee to tell us about things, such as:

mitigating factors that may have impacted or disrupted your studies or exams, and/or whether these have also been considered/factored by the awarding body of your qualifications in determining your grades.`

	cleaned := CleanText(text)

	expected := `In this section
Study
Teacher reference guidance
Study
Imperial Home
StudyApplyUndergraduateApplication processApplication reference
Application reference
Application process
Teacher reference guidance
Everything you need to know about your application reference
UCAS advice: How to get an undergraduate reference
You only need one reference for your UCAS application. It’s usually written for you by a teacher or tutor who knows you – if you’re applying through a school, college or centre registered with UCAS, they’ll do this before they send your completed application to UCAS; you won’t need to do anything for it. See below for more information about choosing a reference if you’re applying independently. What the reference is for
We really value the information provided by your referee as it helps us to gain additional context around your application, in particular about:
your school
your learning environment
the grades you are predicted to achieve
For example, the reference is an opportunity for your referee to tell us about things, such as:
mitigating factors that may have impacted or disrupted your studies or exams, and/or whether these have also been considered/factored by the awarding body of your qualifications in determining your grades.`

	assert.Equal(t, cleaned, expected)
}
