package llm

import (
	"context"
	"reflect"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const (
	openaiModel = openai.GPT4o
)

type OpenAiLLM struct {
	client *openai.Client
}

func NewOpenAiLLM(apiKey string) *OpenAiLLM {
	return &OpenAiLLM{
		client: openai.NewClient(apiKey),
	}
}

func (o *OpenAiLLM) ChatCompletion(ctx context.Context, prompt string) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openaiModel,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	})
	return resp.Choices[0].Message.Content, err
}

func (o *OpenAiLLM) StructuredChatCompletion(ctx context.Context, prompt string, model string) (string, error) {
	return "", nil
}

func openAISchemaFrom(schema interface{}) *jsonschema.Definition {
	s := &jsonschema.Definition{}

	// If schema is a pointer, get the underlying type
	structType := reflect.TypeOf(schema)
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	switch structType.Kind() {
	case reflect.String:
		s.Type = jsonschema.String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s.Type = jsonschema.Integer
	case reflect.Float32, reflect.Float64:
		s.Type = jsonschema.Number
	case reflect.Bool:
		s.Type = jsonschema.Boolean
	case reflect.Struct:
		s.Type = jsonschema.Object
		s.Properties = make(map[string]jsonschema.Definition)

		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			// Split the json tag to handle options like omitempty
			tagParts := strings.Split(jsonTag, ",")
			fieldName := tagParts[0]

			fieldSchema := openAISchemaFrom(reflect.Zero(field.Type).Interface())
			fieldSchema.Description = field.Tag.Get("description")
			s.Properties[fieldName] = *fieldSchema
		}
	case reflect.Slice, reflect.Array:
		s.Type = jsonschema.Array
		s.Items = openAISchemaFrom(reflect.Zero(structType.Elem()).Interface())
	}

	return s
}