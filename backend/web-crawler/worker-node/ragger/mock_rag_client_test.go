package ragger

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMockRagClient(t *testing.T) {
	t.Run("ChunksFrom", func(t *testing.T) {
		mock := NewMockRagClient()
		testInput := "test document"
		expectedChunks := []string{"chunk1", "chunk2"}

		t.Run("returns error when no chunks configured", func(t *testing.T) {
			_, err := mock.ChunksFrom(testInput)
			if err == nil {
				t.Error("expected error when no chunks configured")
			}
		})

		t.Run("returns configured chunks for input", func(t *testing.T) {
			mock.SetChunksFor(testInput, expectedChunks)
			chunks, err := mock.ChunksFrom(testInput)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(chunks, expectedChunks) {
				t.Errorf("got chunks %v, want %v", chunks, expectedChunks)
			}
		})

		t.Run("returns configured error", func(t *testing.T) {
			mock.ChunksError = fmt.Errorf("test error")
			_, err := mock.ChunksFrom(testInput)

			if err == nil {
				t.Error("expected error to be returned")
			}
		})

		t.Run("increments call count", func(t *testing.T) {
			mock := NewMockRagClient()
			mock.SetChunksFor(testInput, expectedChunks)

			initialCount := mock.ChunksCallCount
			_, _ = mock.ChunksFrom(testInput)

			if mock.ChunksCallCount != initialCount+1 {
				t.Errorf("call count not incremented")
			}
		})
	})

	t.Run("ContactsFrom", func(t *testing.T) {
		mock := NewMockRagClient()
		testInput := "test document"
		expectedContacts := []Contact{{
			Value:   "test@example.com",
			Type:    "email",
			Context: "test context",
		}}

		t.Run("returns error when no contacts configured", func(t *testing.T) {
			_, err := mock.ContactsFrom(testInput)
			if err == nil {
				t.Error("expected error when no contacts configured")
			}
		})

		t.Run("returns configured contacts for input", func(t *testing.T) {
			mock.SetContactsFor(testInput, expectedContacts)
			contacts, err := mock.ContactsFrom(testInput)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(contacts, expectedContacts) {
				t.Errorf("got contacts %v, want %v", contacts, expectedContacts)
			}
		})

		t.Run("returns configured error", func(t *testing.T) {
			mock.ContactsError = fmt.Errorf("test error")
			_, err := mock.ContactsFrom(testInput)

			if err == nil {
				t.Error("expected error to be returned")
			}
		})

		t.Run("increments call count", func(t *testing.T) {
			mock := NewMockRagClient()
			mock.SetContactsFor(testInput, expectedContacts)

			initialCount := mock.ContactsCallCount
			_, _ = mock.ContactsFrom(testInput)

			if mock.ContactsCallCount != initialCount+1 {
				t.Errorf("call count not incremented")
			}
		})
	})

	t.Run("EmbeddingsFor", func(t *testing.T) {
		mock := NewMockRagClient()
		testInput := "test document"
		expectedEmbeddings := []float32{0.1, 0.2, 0.3}

		t.Run("returns error when no embeddings configured", func(t *testing.T) {
			_, err := mock.EmbeddingsFor(testInput)
			if err == nil {
				t.Error("expected error when no embeddings configured")
			}
		})

		t.Run("returns configured embeddings for input", func(t *testing.T) {
			mock.SetEmbeddingsFor(testInput, expectedEmbeddings)
			embeddings, err := mock.EmbeddingsFor(testInput)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(embeddings, expectedEmbeddings) {
				t.Errorf("got embeddings %v, want %v", embeddings, expectedEmbeddings)
			}
		})

		t.Run("returns configured error", func(t *testing.T) {
			mock.EmbeddingsError = fmt.Errorf("test error")
			_, err := mock.EmbeddingsFor(testInput)

			if err == nil {
				t.Error("expected error to be returned")
			}
		})

		t.Run("increments call count", func(t *testing.T) {
			mock := NewMockRagClient()
			mock.SetEmbeddingsFor(testInput, expectedEmbeddings)

			initialCount := mock.EmbeddingsCallCount
			_, _ = mock.EmbeddingsFor(testInput)

			if mock.EmbeddingsCallCount != initialCount+1 {
				t.Errorf("call count not incremented")
			}
		})
	})
}

func TestMockRagClient_EmbeddingsForAll(t *testing.T) {
	mock := NewMockRagClient()
	testInput := []string{"test document 1", "test document 2"}
	expectedEmbeddings := [][]float32{{0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}}

	mock.SetEmbeddingsForAll(testInput, expectedEmbeddings)
	embeddings, err := mock.EmbeddingsForAll(testInput)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(embeddings, expectedEmbeddings) {
		t.Errorf("got embeddings %v, want %v", embeddings, expectedEmbeddings)
	}
}
