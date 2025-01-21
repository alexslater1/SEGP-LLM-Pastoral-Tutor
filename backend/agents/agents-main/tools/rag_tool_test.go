package tools

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/knowledge"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestSearchRagFor(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	tool := NewRagTool(knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL")))

	response, err := tool.SearchRagFor("When are the Easter holidays?")

	if err != nil {
		t.Fatalf("Error getting RAG response: %v", err)
	}

	fmt.Println(*response)
}