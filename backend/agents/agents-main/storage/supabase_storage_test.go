package storage

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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

func TestSupabaseStorageUpdate(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	updated, err := Update[Session](storage, "0a54185d-64c3-4a4b-8e98-ccb6cec6c248", map[string]interface{}{"name": "such a good name"})
	if err != nil {
		t.Error("Error updating item in storage", err)
	}

	fmt.Printf("Updated: %+v\n", updated)
}

func TestSupabaseStorageCreateAgentConfig(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	created, err := Store(storage, NewAgentConfig("test", "test", "test", []string{"test"}, []string{"test"}, []string{"test"}))
	if err != nil {
		t.Error("Error creating item in storage", err)
	}

	fmt.Printf("Created: %+v\n", created)
}

func TestSupabaseStorageGetAgentConfig(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := Get[AgentConfig](storage, "a7517604-2f29-47af-855d-81a3808ddf67")
	if err != nil {
		t.Error("Error getting item from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageGetAll(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := storage.getAll(StorageTableNameAgentEvents, NewQueryBuilder().Eq("request_id", "631c06da-6f31-4b5a-9239-62caa40d4c4e"))
	if err != nil {
		t.Error("Error getting items from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageDelete(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := Delete[RequestSession](storage, "5db88b8f-6d92-4d05-8da2-37853d7d1079")
	if err != nil {
		t.Error("Error deleting item from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageDeleteAll(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := DeleteAll[AgentConfig](storage, nil)
	if err != nil {
		t.Error("Error deleting items from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageGetAllChatChecks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := storage.getAll(StorageTableNameChatChecks, nil)
	if err != nil {
		t.Error("Error getting items from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageCreateFeedbackCheck(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := Store(storage, NewFeedbackCheck())
	if err != nil {
		t.Error("Error creating item in storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageGetFeedbackCheck(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := GetAll[FeedbackCheck](storage, NewQueryBuilder().Gt("created_at", time.Now()))
	if err != nil {
		t.Error("Error getting item from storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}
