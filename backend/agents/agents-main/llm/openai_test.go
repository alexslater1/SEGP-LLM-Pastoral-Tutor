package llm

import (
	"reflect"
	"testing"

	"github.com/sashabaranov/go-openai/jsonschema"
)

func TestOpenAISchemaFrom(t *testing.T) {
	type NestedStruct struct {
		NestedField string `json:"nested_field"`
	}

	type TestStruct struct {
		Name        string        `json:"name"`
		Age         int           `json:"age"`
		IsActive    bool          `json:"is_active"`
		Scores      []float64     `json:"scores"`
		Nested      NestedStruct  `json:"nested"`
		Ignored     string        `json:"-"`
		Description string        `json:"description" description:"A description field."`
	}

	tests := []struct {
		name     string
		input    interface{}
		expected *jsonschema.Definition
	}{
		{
			name:  "String type",
			input: "test",
			expected: &jsonschema.Definition{
				Type: jsonschema.String,
			},
		},
		{
			name:  "Integer type",
			input: 42,
			expected: &jsonschema.Definition{
				Type: jsonschema.Integer,
			},
		},
		{
			name:  "Boolean type",
			input: true,
			expected: &jsonschema.Definition{
				Type: jsonschema.Boolean,
			},
		},
		{
			name:  "Float type",
			input: 3.14,
			expected: &jsonschema.Definition{
				Type: jsonschema.Number,
			},
		},
		{
			name:  "Slice type",
			input: []int{1, 2, 3},
			expected: &jsonschema.Definition{
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type: jsonschema.Integer,
				},
			},
		},
		{
			name:  "Struct type",
			input: TestStruct{},
			expected: &jsonschema.Definition{
				Type:       jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"name": {
						Type: jsonschema.String,
					},
					"age": {
						Type: jsonschema.Integer,
					},
					"is_active": {
						Type: jsonschema.Boolean,
					},
					"scores": {
						Type:  jsonschema.Array,
						Items: &jsonschema.Definition{Type: jsonschema.Number},
					},
					"nested": {
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"nested_field": {
								Type: jsonschema.String,
							},
						},
					},
					"description": {
						Type:        jsonschema.String,
						Description: "A description field.",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := openAISchemaFrom(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("unexpected result: got %+v, want %+v", result, test.expected)
			}
		})
	}
}
