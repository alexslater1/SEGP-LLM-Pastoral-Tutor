package llm

import (
	"testing"

	"github.com/google/generative-ai-go/genai"
	"github.com/stretchr/testify/assert"
)

func TestGenaiSchemaFrom(t *testing.T) {
	t.Run("primitive types", func(t *testing.T) {
		tests := []struct {
			name     string
			input    interface{}
			expected *genai.Schema
		}{
			{
				name:  "string",
				input: "test",
				expected: &genai.Schema{
					Type: genai.TypeString,
				},
			},
			{
				name:  "int",
				input: 42,
				expected: &genai.Schema{
					Type: genai.TypeInteger,
				},
			},
			{
				name:  "int32",
				input: int32(42),
				expected: &genai.Schema{
					Type:   genai.TypeInteger,
					Format: "int32",
				},
			},
			{
				name:  "int64",
				input: int64(42),
				expected: &genai.Schema{
					Type:   genai.TypeInteger,
					Format: "int64",
				},
			},
			{
				name:  "float32",
				input: float32(3.14),
				expected: &genai.Schema{
					Type:   genai.TypeNumber,
					Format: "float",
				},
			},
			{
				name:  "float64",
				input: 3.14,
				expected: &genai.Schema{
					Type:   genai.TypeNumber,
					Format: "double",
				},
			},
			{
				name:  "bool",
				input: true,
				expected: &genai.Schema{
					Type: genai.TypeBoolean,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := genaiSchemaFrom(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("struct types", func(t *testing.T) {
		type TestStruct struct {
			Name        string  `json:"name" description:"The name field"`
			Age         int     `json:"age"`
			Score       float64 `json:"score"`
			IsActive    bool    `json:"is_active"`
			IgnoreField string  `json:"-"`
		}

		input := TestStruct{}
		result := genaiSchemaFrom(input)

		expected := &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"name": {
					Type:        genai.TypeString,
					Description: "The name field",
				},
				"age": {
					Type: genai.TypeInteger,
				},
				"score": {
					Type:   genai.TypeNumber,
					Format: "double",
				},
				"is_active": {
					Type: genai.TypeBoolean,
				},
			},
			Required: []string{"name", "age", "score", "is_active"},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("slice types", func(t *testing.T) {
		input := []string{}
		result := genaiSchemaFrom(input)

		expected := &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeString,
			},
		}

		assert.Equal(t, expected, result)
	})

	t.Run("pointer types", func(t *testing.T) {
		str := "test"
		result := genaiSchemaFrom(&str)

		expected := &genai.Schema{
			Type: genai.TypeString,
		}

		assert.Equal(t, expected, result)
	})
}

func TestGenaiSchemaFromComplexStructure(t *testing.T) {
	type Address struct {
		Street  string   `json:"street"`
		City    string   `json:"city"`
		ZipCode int      `json:"zip_code"`
		Tags    []string `json:"tags" description:"Location tags"`
	}

	type Contact struct {
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Preferred []string `json:"preferred_contact_methods"`
	}

	type Department struct {
		Name     string    `json:"name"`
		Location *Address  `json:"location"`
		Teams    []string  `json:"teams"`
		Budget   float64   `json:"budget"`
		Active   bool      `json:"active"`
		Leads    []Contact `json:"leads"`
	}

	type Company struct {
		Name        string       `json:"name" description:"Company's legal name"`
		Founded     int32        `json:"founded"`
		Departments []Department `json:"departments"`
		HeadOffice  Address      `json:"head_office"`
		MainContact Contact      `json:"main_contact"`
	}

	input := Company{}
	result := genaiSchemaFrom(input)

	expected := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"name": {
				Type:        genai.TypeString,
				Description: "Company's legal name",
			},
			"founded": {
				Type:   genai.TypeInteger,
				Format: "int32",
			},
			"departments": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name": {
							Type: genai.TypeString,
						},
						"location": {
							Type: genai.TypeObject,
							Properties: map[string]*genai.Schema{
								"street":   {Type: genai.TypeString},
								"city":     {Type: genai.TypeString},
								"zip_code": {Type: genai.TypeInteger},
								"tags": {
									Type:        genai.TypeArray,
									Items:       &genai.Schema{Type: genai.TypeString},
									Description: "Location tags",
								},
							},
							Required: []string{"street", "city", "zip_code", "tags"},
						},
						"teams": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
						"budget": {
							Type:   genai.TypeNumber,
							Format: "double",
						},
						"active": {
							Type: genai.TypeBoolean,
						},
						"leads": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"email": {Type: genai.TypeString},
									"phone": {Type: genai.TypeString},
									"preferred_contact_methods": {
										Type:  genai.TypeArray,
										Items: &genai.Schema{Type: genai.TypeString},
									},
								},
								Required: []string{"email", "phone", "preferred_contact_methods"},
							},
						},
					},
					Required: []string{"name", "location", "teams", "budget", "active", "leads"},
				},
			},
			"head_office": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"street":   {Type: genai.TypeString},
					"city":     {Type: genai.TypeString},
					"zip_code": {Type: genai.TypeInteger},
					"tags": {
						Type:        genai.TypeArray,
						Items:       &genai.Schema{Type: genai.TypeString},
						Description: "Location tags",
					},
				},
				Required: []string{"street", "city", "zip_code", "tags"},
			},
			"main_contact": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"email": {Type: genai.TypeString},
					"phone": {Type: genai.TypeString},
					"preferred_contact_methods": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
				},
				Required: []string{"email", "phone", "preferred_contact_methods"},
			},
		},
		Required: []string{"name", "founded", "departments", "head_office", "main_contact"},
	}

	assert.Equal(t, expected, result)
}

func TestGenaiSchemaFromSliceOfComplex(t *testing.T) {
	type Address struct {
		Street  string   `json:"street"`
		City    string   `json:"city"`
		ZipCode int      `json:"zip_code"`
		Tags    []string `json:"tags" description:"Location tags"`
	}

	type Contact struct {
		Email     string   `json:"email"`
		Phone     string   `json:"phone"`
		Preferred []string `json:"preferred_contact_methods"`
	}

	type Department struct {
		Name     string    `json:"name"`
		Location *Address  `json:"location"`
		Teams    []string  `json:"teams"`
		Budget   float64   `json:"budget"`
		Active   bool      `json:"active"`
		Leads    []Contact `json:"leads"`
	}

	type Company struct {
		Name        string       `json:"name" description:"Company's legal name"`
		Founded     int32        `json:"founded"`
		Departments []Department `json:"departments"`
		HeadOffice  Address      `json:"head_office"`
		MainContact Contact      `json:"main_contact"`
	}

	input := []Company{}
	result := genaiSchemaFrom(input)

	expected := &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"name": {
					Type:        genai.TypeString,
					Description: "Company's legal name",
				},
				"founded": {
					Type:   genai.TypeInteger,
					Format: "int32",
				},
				"departments": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"name": {
								Type: genai.TypeString,
							},
							"location": {
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"street":   {Type: genai.TypeString},
									"city":     {Type: genai.TypeString},
									"zip_code": {Type: genai.TypeInteger},
									"tags": {
										Type:        genai.TypeArray,
										Items:       &genai.Schema{Type: genai.TypeString},
										Description: "Location tags",
									},
								},
								Required: []string{"street", "city", "zip_code", "tags"},
							},
							"teams": {
								Type:  genai.TypeArray,
								Items: &genai.Schema{Type: genai.TypeString},
							},
							"budget": {
								Type:   genai.TypeNumber,
								Format: "double",
							},
							"active": {
								Type: genai.TypeBoolean,
							},
							"leads": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type: genai.TypeObject,
									Properties: map[string]*genai.Schema{
										"email": {Type: genai.TypeString},
										"phone": {Type: genai.TypeString},
										"preferred_contact_methods": {
											Type:  genai.TypeArray,
											Items: &genai.Schema{Type: genai.TypeString},
										},
									},
									Required: []string{"email", "phone", "preferred_contact_methods"},
								},
							},
						},
						Required: []string{"name", "location", "teams", "budget", "active", "leads"},
					},
				},
				"head_office": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"street":   {Type: genai.TypeString},
						"city":     {Type: genai.TypeString},
						"zip_code": {Type: genai.TypeInteger},
						"tags": {
							Type:        genai.TypeArray,
							Items:       &genai.Schema{Type: genai.TypeString},
							Description: "Location tags",
						},
					},
					Required: []string{"street", "city", "zip_code", "tags"},
				},
				"main_contact": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"email": {Type: genai.TypeString},
						"phone": {Type: genai.TypeString},
						"preferred_contact_methods": {
							Type:  genai.TypeArray,
							Items: &genai.Schema{Type: genai.TypeString},
						},
					},
					Required: []string{"email", "phone", "preferred_contact_methods"},
				},
			},
			Required: []string{"name", "founded", "departments", "head_office", "main_contact"},
		},
	}

	assert.Equal(t, expected, result)
}
