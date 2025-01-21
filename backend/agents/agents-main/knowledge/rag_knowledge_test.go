package knowledge

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestRAGKnowledge(t *testing.T) {
	rag := NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
	result, err := rag.Get("When is the exam period?s")
	if err != nil {
		t.Fatalf("Error getting RAG response: %v", err)
	}
	t.Logf("RAG response: %s", *result)
}
