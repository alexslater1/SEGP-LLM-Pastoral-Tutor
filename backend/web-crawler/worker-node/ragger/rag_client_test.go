package ragger

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactsFrom(t *testing.T) {
	client := &RAGClient{}

	text := `You can contact me at john.doe@example.com or call me at (555) 123-4567.
My website is https://www.example.com.
Alternatively, reach out to jane_doe123@example.org for further details.`

	contacts, err := client.ContactsFrom(text)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Define the expected contacts.
	// Note: The extraction order is determined by the order of regex processing:
	// 1. Emails (in order of appearance)
	// 2. Phone numbers
	// 3. Addresses
	expected := []struct {
		value string
		ctype ContactType
	}{
		{"john.doe@example.com", ContactTypeEmail},
		{"jane_doe123@example.org", ContactTypeEmail},
		{"(555) 123-4567", ContactTypePhone},
		{"https://www.example.com", ContactTypeWebsite},
	}

	// Check that the number of contacts is as expected.
	if len(contacts) != len(expected) {
		t.Fatalf("Expected %d contacts, got %d", len(expected), len(contacts))
	}

	// Validate each contact.
	for i, exp := range expected {
		got := contacts[i]

		// Check the type.
		if got.Type != exp.ctype {
			t.Errorf("Contact %d: expected type %q, got %q", i, exp.ctype, got.Type)
		}

		// Check the value.
		if got.Value != exp.value {
			t.Errorf("Contact %d: expected value %q, got %q", i, exp.value, got.Value)
		}

		// Check that the context contains the value.
		if !strings.Contains(got.Context, got.Value) {
			t.Errorf("Contact %d: context %q does not contain the value %q", i, got.Context, got.Value)
		}
	}
}

func TestContactsFrom_NoContacts(t *testing.T) {
	client := &RAGClient{}

	text := "Hello there my name is Ethan"

	contacts, err := client.ContactsFrom(text)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(contacts))
}

func TestChunksFrom(t *testing.T) {
	client := NewRAGClient(modelPath, libraryPath, tokenizerPath)

	text := "Hello there my name is Ethan"

	chunks, err := client.ChunksFrom(text)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(chunks))
}

func TestEmbeddingsFor(t *testing.T) {
	client := NewRAGClient(modelPath, libraryPath, tokenizerPath)

	text := "Hello there my name is Ethan"

	embeddings, err := client.EmbeddingsFor(text)
	assert.NoError(t, err)
	assert.Equal(t, 384, len(embeddings))
}

func TestEmbeddingsForAll(t *testing.T) {
	client := NewRAGClient(modelPath, libraryPath, tokenizerPath)

	texts := []string{"Hello there my name is Ethan", "Hello there my name is Ethan"}

	embeddings, err := client.EmbeddingsForAll(texts)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(embeddings))
	assert.Equal(t, 384, len(embeddings[0]))
	assert.Equal(t, 384, len(embeddings[1]))
}
