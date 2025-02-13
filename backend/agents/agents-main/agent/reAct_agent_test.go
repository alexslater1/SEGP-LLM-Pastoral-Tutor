package agent

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/context_keys"
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

	ctx := context.Background()
	ctx = context_keys.SetRequestID(ctx, "test-request-id")

	response, reasoning, err := agent.Run(ctx, "What is the current price of the dollar")
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

	createdReq, err := storage.Store(supabaseStorage, storage.NewAgentRequest("/test", map[string]string{"test": "test"}, ""))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	ctx := context_keys.SetRequestID(context.Background(), createdReq.ID)

	response, reasoning, err := agent.Run(ctx, "What is the current price of the dollar")
	if err != nil {
		t.Fatalf("Error running agent: %v", err)
	}

	t.Logf("Response: %s", *response)
	t.Logf("Reasoning: %s", *reasoning)
}
