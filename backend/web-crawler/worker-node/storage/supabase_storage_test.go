package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethanhosier/worker-node/ragger"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	modelPath     = filepath.Join("..", "model", "model.onnx")
	libraryPath   = filepath.Join("..", "libonnxruntime.so.1.20.1")
	tokenizerPath = filepath.Join("..", "model", "tokenizer.json")
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestParseStringToFloat32Slice(t *testing.T) {
	var embedding []float32
	err := parseStringToFloat32Slice("[1.5,2.5,3.5]", &embedding)
	if err != nil {
		t.Error("Error parsing string to float32 slice", err)
	}

	assert.Equal(t, []float32{1.5, 2.5, 3.5}, embedding)
}

func TestProcessEmbeddingField(t *testing.T) {
	data := map[string]interface{}{
		"embedding": "[1.5,2.5,3.5]",
	}

	processed, err := processEmbeddingField(data)
	if err != nil {
		t.Error("Error processing embedding field", err)
	}

	processedMap := processed.(map[string]interface{})
	assert.Equal(t, []float32{1.5, 2.5, 3.5}, processedMap["embedding"])
}

func TestProcessEmbeddingFieldNoEmbedding(t *testing.T) {
	data := map[string]interface{}{
		"test": "test",
	}

	processed, err := processEmbeddingField(data)
	if err != nil {
		t.Error("Error processing embedding field", err)
	}

	assert.Equal(t, data, processed)
}

func TestProcessAllEmbeddingFields(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"embedding": "[1.5,2.5,3.5]"},
		map[string]interface{}{"embedding": "[4.5,5.5,6.5]"},
	}

	processed, err := processAllEmbeddingFields(data)
	if err != nil {
		t.Error("Error processing all embedding fields", err)
	}

	assert.Equal(t, []interface{}{
		map[string]interface{}{"embedding": []float32{1.5, 2.5, 3.5}},
		map[string]interface{}{"embedding": []float32{4.5, 5.5, 6.5}},
	}, processed)
}

func TestSupabaseStorageStore(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	storage := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	data, err := storage.store(StorageTableNameAgentRequests, NewAgentRequest("test", map[string]string{"test": "test"}))
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

	data, err := Store(storage, NewAgentRequest("test", map[string]string{"test": "test"}))
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

	data, err := storage.get(StorageTableNameAgentRequests, "00c4a620-d375-49ab-b5bd-0b67d4755fd8")
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

	data, err := Store(storage, NewAgentEvent("00c4a620-d375-49ab-b5bd-0b67d4755fd8", "test", map[string]string{"test": "test"}))
	if err != nil {
		t.Error("Error storing item in storage", err)
	}

	fmt.Printf("Data: %+v\n", data)
}

func TestSupabaseStorageTestStoreRagSource(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	store := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	ragSource := RagSource{
		URL: "https://example.com",
	}

	_, err := Store(store, ragSource)
	if err != nil {
		t.Error("Error storing item in storage", err)
	}
}

func TestSupabaseStorageTestStoreRagChunk(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	store := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	ragger := ragger.NewRAGClient(modelPath, libraryPath, tokenizerPath)
	embedding, err := ragger.EmbeddingsFor("Hello, world!")
	if err != nil {
		t.Error("Error embedding text", err)
	}

	ragChunk := RagChunk{
		Text:        "Hello, world!",
		PosInSource: 1,
		Embedding:   embedding,
		RagSourceId: 75,
	}

	storedRagChunk, err := Store(store, ragChunk)
	if err != nil {
		t.Error("Error storing item in storage", err)
	}

	assert.Equal(t, ragChunk.Embedding, storedRagChunk.Embedding)
	assert.Equal(t, ragChunk.Text, storedRagChunk.Text)
	assert.Equal(t, ragChunk.PosInSource, storedRagChunk.PosInSource)
	assert.Equal(t, ragChunk.RagSourceId, storedRagChunk.RagSourceId)
}

func TestSupabaseStorageTestStoreRagContacts(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	store := NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	ragger := ragger.NewRAGClient(modelPath, libraryPath, tokenizerPath)
	embedding, err := ragger.EmbeddingsFor("Hello, world!")
	if err != nil {
		t.Error("Error embedding text", err)
	}

	ragContact := RagContact{
		RagSourceId: 75,
		Context:     "test",
		Contact:     "test@example.com",
		PosInSource: 1,
		ContactType: "test",
		Embedding:   embedding,
	}

	storedRagContact, err := Store(store, ragContact)
	if err != nil {
		t.Error("Error storing item in storage", err)
	}

	assert.Equal(t, ragContact.Context, storedRagContact.Context)
	assert.Equal(t, ragContact.Contact, storedRagContact.Contact)
	assert.Equal(t, ragContact.PosInSource, storedRagContact.PosInSource)
	assert.Equal(t, ragContact.ContactType, storedRagContact.ContactType)
	assert.Equal(t, ragContact.RagSourceId, storedRagContact.RagSourceId)
	assert.Equal(t, ragContact.Embedding, storedRagContact.Embedding)
}
