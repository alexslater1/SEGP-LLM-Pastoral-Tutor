package llm

import (
	"context"
	"reflect"
	"strings"

	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiLLM struct {
	client *genai.Client
}

func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiLLM, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &GeminiLLM{
		client: client,
	}, nil
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
