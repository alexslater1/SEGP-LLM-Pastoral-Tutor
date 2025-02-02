package llm

import (
	"context"
	"os"
	"testing"

	"github.com/segp/agents-main/tools"
)

func TestDeepSeekLLM(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD environment")
	}

	ocClient := NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))

	prompt := "What is the capital of France?"

	response, err := ocClient.ChatCompletion(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestDeepSeekCompletionthinking(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD environment")
	}

	ocClient := NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))

	prompt := "You are given a list of tools. Pick the best tool for this query `what is the current price of the usd`. Tools: [`no_tool`:`pick no tool, either as know the answer or no relevant tool`, `google search`: `get the contents of the top 3 search results for a query`, `google_maps`: `get the contents of the top 3 search results for a query`] Include thinking in <thoughts> </thoughts> tags and then answer in <answer> </answer> tags. Answer must include the tool name and the tool input."

	response, err := ocClient.ChatCompletion(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error calling ChatCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestDeepSeekStructuredOutput(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD environment")
	}

	ocClient := NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))

	prompt := "Give me a random address"

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		State  string `json:"state"`
		Zip    string `json:"zip"`
	}

	response, err := ocClient.StructuredOutputCompletion(context.Background(), prompt, Address{})
	if err != nil {
		t.Fatalf("Error calling StructuredOutputCompletion: %v", err)
	}

	t.Logf("Response: %v", *response)
}

func TestDeepSeekChatCompletionWithTools(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD environment")
	}

	ocClient := NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))

	prompt := "What is the weather in Paris on 2025-01-14?"

	toolCalls, err := ocClient.ChatCompletionWithTools(context.Background(), prompt, []tools.ToolDefinition{tools.CheckWeatherTool(), tools.CheckBestAnimalNameTool()}, tools.ToolChoice{Type: tools.ToolChoiceTypeAuto})
	if err != nil {
		t.Fatalf("Error calling ChatCompletionWithTools: %v", err)
	}

	t.Logf("Tool calls: %v", toolCalls)
}
