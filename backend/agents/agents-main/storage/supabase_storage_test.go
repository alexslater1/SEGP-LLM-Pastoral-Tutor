package storage

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestSupabaseStorageStore(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	storage.store(StorageTableNameAgentRequests, AgentRequest{
		Endpoint: "test",
	})
}

func TestSupabaseStorageGet(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := storage.get(StorageTableNameAgentRequests, "dd774b50-4843-4615-9529-160a65391ae8")
	if err != nil {
		t.Error("Error getting item from storage")
	}

	fmt.Printf("%+v\n", data)
}
