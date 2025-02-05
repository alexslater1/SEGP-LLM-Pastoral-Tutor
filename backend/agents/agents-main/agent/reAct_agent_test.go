package agent

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/storage"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestReActAgent(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	agent := NewDefaultLoggingReActAgent()

	response, reasoning, err := agent.Run("What is the current price of the dollar", "test-request-id")
	if err != nil {
		t.Fatalf("Error running agent: %v", err)
	}

	t.Logf("Response: %s", *response)
	t.Logf("Reasoning: %s", *reasoning)
}

func TestEventStoringReActAgent(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}
	
	agent := NewDefaultEventStoringLoggingReActAgent()
	supabaseStorage := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	createdReq, err := storage.Store(supabaseStorage, storage.NewAgentRequest("/test", map[string]string{"test": "test"}))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	response, reasoning, err := agent.Run("What is the current price of the dollar", createdReq.ID)
	if err != nil {
		t.Fatalf("Error running agent: %v", err)
	}

	t.Logf("Response: %s", *response)
	t.Logf("Reasoning: %s", *reasoning)
}
