package llm

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/generative-ai-go/genai"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/tools"
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

func TestGemini(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	prompt := "What is the capital of France?"

	response, err := geminiClient.ChatCompletion(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestGeminiStructuredOutput(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	prompt := "Give me a random address"

	type Address struct {
		Street  string   `json:"street"`
		City    string   `json:"city"`
		ZipCode int      `json:"zip_code"`
		Tags    []string `json:"tags" description:"Location tags"`
	}

	response, err := geminiClient.StructuredOutputCompletion(context.Background(), prompt, Address{})
	if err != nil {
		t.Fatalf("Error calling StructuredOutputCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestGeminiChatCompletionWithTools(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	prompt := "What is the weather in San Francisco on 10/10/2024 and in New York on 10/10/2024?"

	tool := tools.CheckWeatherTool()

	response, err := geminiClient.ChatCompletionWithTools(context.Background(), prompt, []tools.ToolDefinition{tool}, tools.ToolChoice{Type: tools.ToolChoiceTypeAuto})
	if err != nil {
		t.Fatalf("Error calling ChatCompletionWithTools: %v", err)
	}

	t.Logf("Response: %+v", response)
}

func TestGeminiChatCompletionLLMThinking(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	prompt := "You are given a list of tools. Pick the best tool for this query `what is the current price of the usd`. Tools: [`no_tool`:`pick no tool, either as know the answer or no relevant tool`, `google search`: `get the contents of the top 3 search results for a query`, `google_maps`: `get the contents of the top 3 search results for a query`] Include thinking in <thoughts> </thoughts> tags and then answer in <answer> </answer> tags. Answer must include the tool name and the tool input."

	response, err := geminiClient.ChatCompletion(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestGeminiChatCompletionLLMThinkingWithStructuredOutput(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	prompt := "You are to answer the query `what is the weather in the home country of yesterday's new richest man`. It is essential that you first think about all steps required to solve this problem, and must pick the next ONE tool which should be called to help next. Tools: [`no_tool`:`pick no tool, either as know the answer or no relevant tool`, `google search`: `get the contents of the top 3 search results for a query`, `google_maps`: `get the contents of the top 3 search results for a query`]"

	type ToolParam struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type Response struct {
		Thoughts string `json:"_thoughts"`
		Answer   struct {
			ToolName   string      `json:"tool_name"`
			ToolParams []ToolParam `json:"tool_params"`
		} `json:"answer"`
	}

	ts := []tools.ToolDefinition{tools.NewGoogleSearchResultsTool(nil).Definition(), tools.NewGiveAnswerTool().Definition()}
	definitionStr := ""

	for _, t := range ts {
		definitionStr += fmt.Sprintf("%s: %s\n", t.Name, t.Description)
		for _, p := range t.Parameters {
			definitionStr += fmt.Sprintf("%s: %s\n", p.Name, p.Description)
		}
	}

	response, err := geminiClient.StructuredOutputCompletion(context.Background(), prompt, Response{})
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %+v", *response)
}

func TestStructuredOutputCompletionWithTools(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	geminiClient := NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY"))

	gs := googleSearch.NewRodClient()

	ts := []tools.ToolDefinition{tools.NewGoogleSearchResultsTool(gs).Definition(), tools.NewGiveAnswerTool().Definition()}

	prompt := "What is the current price of the usd"

	response, err := geminiClient.ChatCompletionWithTools(context.Background(), prompt, ts, tools.ToolChoice{Type: tools.ToolChoiceTypeRequired})
	if err != nil {
		t.Fatalf("Error calling ChatCompletionWithTools: %v", err)
	}

	t.Logf("Response: %+v", response[0])
}
