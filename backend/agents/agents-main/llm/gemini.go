package llm

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/segp/agents-main/tools"
	"google.golang.org/api/option"
)

const (
	chatCompletionModel = "gemini-2.0-flash"
)

type GeminiLLM struct {
	client *genai.Client
}

func NewGeminiLLM(ctx context.Context, apiKey string) *GeminiLLM {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	return &GeminiLLM{
		client: client,
	}
}

func (g *GeminiLLM) ChatCompletion(ctx context.Context, prompt string) (*string, error) {
	model := g.client.GenerativeModel(chatCompletionModel)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates received")
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response parts received")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("failed to convert to text")
	}

	strText := string(text)

	return &strText, nil
}

func (g *GeminiLLM) StructuredOutputCompletion(ctx context.Context, prompt string, schema interface{}) (*string, error) {
	model := g.client.GenerativeModel(chatCompletionModel)
	responseSchema := genaiSchemaFrom(schema)

	if responseSchema != nil {
		model.ResponseMIMEType = "application/json"
		model.ResponseSchema = responseSchema
	}

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no response candidates received")
	}

	text, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("failed to convert content to text")
	}

	strText := string(text)
	return &strText, nil
}

func (g *GeminiLLM) ChatCompletionWithTools(ctx context.Context, prompt string, ts []tools.ToolDefinition, toolChoice tools.ToolChoice) ([]tools.ToolCall, error) {
	model := g.client.GenerativeModel(chatCompletionModel)

	genaiTools := []*genai.Tool{}
	for _, tool := range ts {
		genaiTools = append(genaiTools, genaiToolFrom(tool))
	}

	genaiToolConfig := genai.ToolConfig{
		FunctionCallingConfig: functionCallingConfigFrom(toolChoice),
	}

	model.Tools = genaiTools
	model.ToolConfig = &genaiToolConfig

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response candidates received")
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response parts received")
	}

	toolCalls := []tools.ToolCall{}
	for _, part := range resp.Candidates[0].Content.Parts {
		switch p := part.(type) {
		case genai.FunctionCall:
			toolCalls = append(toolCalls, toolCallFromGenai(p))
		case genai.Text:
		default:
			fmt.Printf("Unknown type: %T\n", p)
		}
	}

	return toolCalls, nil
}

func genaiSchemaFrom(schema interface{}) *genai.Schema {
	s := &genai.Schema{}

	// If schema is a pointer, get the underlying type
	structType := reflect.TypeOf(schema)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	switch structType.Kind() {
	case reflect.String:
		s.Type = genai.TypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = genai.TypeInteger
		// Optionally set format based on size
		if structType.Kind() == reflect.Int32 {
			s.Format = "int32"
		} else if structType.Kind() == reflect.Int64 {
			s.Format = "int64"
		}
	case reflect.Float32, reflect.Float64:
		s.Type = genai.TypeNumber
		// Set format based on size
		if structType.Kind() == reflect.Float32 {
			s.Format = "float"
		} else {
			s.Format = "double"
		}
	case reflect.Bool:
		s.Type = genai.TypeBoolean
	case reflect.Struct:
		s.Type = genai.TypeObject
		s.Properties = make(map[string]*genai.Schema)

		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			s.Required = append(s.Required, jsonTag)

			// Split the json tag to handle options like omitempty
			tagParts := strings.Split(jsonTag, ",")
			fieldName := tagParts[0]

			s.Properties[fieldName] = genaiSchemaFrom(reflect.Zero(field.Type).Interface())
			// You might want to add field.Tag.Get("description") to the schema description
			s.Properties[fieldName].Description = field.Tag.Get("description")
		}
	case reflect.Slice, reflect.Array:
		s.Type = genai.TypeArray
		s.Items = genaiSchemaFrom(reflect.Zero(structType.Elem()).Interface())
	}

	return s
}
