package agent

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestReActAgent(t *testing.T) {
	agent := NewDefaultReActAgent()

	response, reasoning, err := agent.Run("Who is the headmaster of the secondary school which Dillan Scott (imperial college london) attended?")
	if err != nil {
		t.Fatalf("Error running agent: %v", err)
	}

	t.Logf("Response: %s", *response)
	t.Logf("Reasoning: %s", *reasoning)
}
