package llm

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env file before running tests
	err := godotenv.Load("../../.env")
	if err != nil {
		// Don't fail if .env file is not found, just log it
		println("Warning: .env file not found")
	}

	os.Exit(m.Run())
}

func TestOpenAIChatCompletion(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	ocClient := NewOpenAiLLM(os.Getenv("OPENAI_API_KEY"))

	prompt := "What is the capital of France?"

	response, err := ocClient.ChatCompletion(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestOpenAIChatCompletionWithSchema(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD environment")
	}

	ocClient := NewOpenAiLLM(os.Getenv("OPENAI_API_KEY"))

	prompt := "Give me a random address"

	type Address struct {
		Street  string   `json:"street"`
		City    string   `json:"city"`
		ZipCode int      `json:"zip_code"`
		Tags    []string `json:"tags" description:"Location tags"`
	}

	response, err := ocClient.StructuredOutputCompletion(context.Background(), prompt, Address{})
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}
