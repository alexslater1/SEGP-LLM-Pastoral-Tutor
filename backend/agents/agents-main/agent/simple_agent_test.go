package agent

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/segp/agents-main/llm"
	"github.com/stretchr/testify/assert"
)

func TestSimpleAgent(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	var (
		systemPrompt = "You are a chat title creator agent. You will receive a query and you will need to create a title for the chat. The title should be a short sentence which captures the essence of the query. Return with just the title"
		userPrompt   = "What is the weather in Tokyo"
	)

	agent := NewSimpleAgent(systemPrompt, llm.NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY")))
	assert.NotNil(t, agent)

	answer, reason, err := agent.Run(userPrompt, "123")
	assert.NoError(t, err)
	assert.NotNil(t, answer)
	assert.NotNil(t, reason)

	fmt.Println("Answer:", *answer)
	fmt.Println("Reason:", *reason)
}
