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

	data, err := storage.store(StorageTableNameAgentRequests, NewAgentRequest("test", map[string]string{"test": "test"}, ""))
	if err != nil {
		t.Error("Error storing item in storage")
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestStoreSupabaseStorageStore(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := Store(storage, NewAgentRequest("test", map[string]string{"test": "test"}, ""))
	if err != nil {
		t.Error("Error storing item in storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageGet(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := storage.get(StorageTableNameAgentRequests, "1e69f039-a124-41c9-8506-8615164cca4b")
	if err != nil {
		t.Error("Error getting item from storage")
	}

	fmt.Printf("%+v\n", data)
}

func TestStoreSupabaseStorageStoreAgentEvent(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := Store(storage, NewAgentEvent("e3f3e50b-bd95-46f4-974f-e1ba6f91fa5f", "test", map[string]string{"test": "test"}))
	if err != nil {
		t.Error("Error storing item in storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}
