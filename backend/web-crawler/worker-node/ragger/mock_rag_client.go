package ragger

import (
	"fmt"
	"strings"
)

// MockRagClient implements the Ragger interface for testing purposes
type MockRagClient struct {
	// Maps input text to return values
	ChunksMap           map[string][]string
	ContactsMap         map[string][]Contact
	EmbeddingsMap       map[string][]float32
	EmbeddingsForAllMap map[string][][]float32
	// Error states
	ChunksError           error
	ContactsError         error
	EmbeddingsError       error
	EmbeddingsForAllError error
	// Call counters for verification
	ChunksCallCount           int
	ContactsCallCount         int
	EmbeddingsCallCount       int
	EmbeddingsForAllCallCount int
}

// NewMockRagClient creates a new instance of MockRagClient
func NewMockRagClient() *MockRagClient {
	return &MockRagClient{
		ChunksMap:           make(map[string][]string),
		ContactsMap:         make(map[string][]Contact),
		EmbeddingsMap:       make(map[string][]float32),
		EmbeddingsForAllMap: make(map[string][][]float32),
	}
}

// SetChunksFor sets the chunks to return for a specific input text
func (m *MockRagClient) SetChunksFor(input string, chunks []string) {
	m.ChunksMap[input] = chunks
}

// SetContactsFor sets the contacts to return for a specific input text
func (m *MockRagClient) SetContactsFor(input string, contacts []Contact) {
	m.ContactsMap[input] = contacts
}

// SetEmbeddingsFor sets the embeddings to return for a specific input text
func (m *MockRagClient) SetEmbeddingsFor(input string, embeddings []float32) {
	m.EmbeddingsMap[input] = embeddings
}

func (m *MockRagClient) SetEmbeddingsForAll(input []string, embeddings [][]float32) {
	m.EmbeddingsForAllMap[embeddingsForAllKey(input)] = embeddings
}

func (m *MockRagClient) ChunksFrom(text string) ([]string, error) {
	m.ChunksCallCount++
	if m.ChunksError != nil {
		return nil, m.ChunksError
	}
	chunks, ok := m.ChunksMap[text]
	if !ok {
		return nil, fmt.Errorf("no chunks configured for input: %s", text)
	}
	return chunks, nil
}

func (m *MockRagClient) ContactsFrom(text string) ([]Contact, error) {
	m.ContactsCallCount++
	if m.ContactsError != nil {
		return nil, m.ContactsError
	}
	contacts, ok := m.ContactsMap[text]
	if !ok {
		return nil, fmt.Errorf("no contacts configured for input: %s", text)
	}
	return contacts, nil
}

func (m *MockRagClient) EmbeddingsFor(text string) ([]float32, error) {
	m.EmbeddingsCallCount++
	if m.EmbeddingsError != nil {
		return nil, m.EmbeddingsError
	}
	embeddings, ok := m.EmbeddingsMap[text]
	if !ok {
		return nil, fmt.Errorf("no embeddings configured for input: %s", text)
	}
	return embeddings, nil
}

func (m *MockRagClient) EmbeddingsForAll(texts []string) ([][]float32, error) {
	m.EmbeddingsForAllCallCount++
	if m.EmbeddingsForAllError != nil {
		return nil, m.EmbeddingsForAllError
	}
	embeddings, ok := m.EmbeddingsForAllMap[embeddingsForAllKey(texts)]
	if !ok {
		return nil, fmt.Errorf("no embeddings configured for input: %s", texts)
	}
	return embeddings, nil
}

func embeddingsForAllKey(texts []string) string {
	return strings.Join(texts, "####")
}
