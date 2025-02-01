package ragger

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

const (
	contactContextLengthCharsBefore = 200
	contactContextLengthCharsAfter  = 50
)

type RAGClient struct {
	embedder  *Embedder
	chunker   *Chunker
	tokenizer *tokenizer.Tokenizer
}

func NewRAGClient(modelPath string, libraryPath string, tokenizerPath string) *RAGClient {
	tok, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		log.Fatalf("Error loading tokenizer: %v", err)
	}

	embedder, err := NewEmbedder(modelPath, libraryPath, tok)
	if err != nil {
		log.Fatalf("Error loading embedder: %v", err)
	}

	return &RAGClient{
		embedder:  embedder,
		chunker:   NewChunker(tok),
		tokenizer: tok,
	}
}

func (c *RAGClient) ChunksFrom(text string) ([]string, error) {
	return c.chunker.Chunk(text)
}

func (c *RAGClient) EmbeddingsFor(text string) ([]float32, error) {
	return c.embedder.Embed(text)
}

func (c *RAGClient) EmbeddingsForAll(texts []string) ([][]float32, error) {
	return c.embedder.EmbedAll(texts)
}

func (c *RAGClient) ContactsFrom(text string) ([]Contact, error) {
	var result []Contact
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from in contactsFrom panic: %v\n", r)
			// Set result to empty contacts on panic
			result = []Contact{}
		}
	}()

	var contacts []Contact

	// Regex for emails.
	emailRegex := regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	// Regex for phone numbers.
	// This pattern expects the phone number to be in the format (xxx) xxx-xxxx.
	phoneRegex := regexp.MustCompile(`\(\d{3}\)[\s-]*\d{3}[\s-]*\d{4}`)
	// Regex for websites.
	// Note: This simple pattern may include trailing punctuation.
	websiteRegex := regexp.MustCompile(`(?i)\b(?:https?://|www\.)[^\s]+`)

	// Helper function to compute context: 200 chars before, 50 chars after.
	getContext := func(start, end int) string {
		from := start - 200
		if from < 0 {
			from = 0
		}
		to := end + 50
		if to > len(text) {
			to = len(text)
		}
		context := text[from:to]

		// Check token count - return empty string if too long
		enc, err := c.tokenizer.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(context)), true)
		if err != nil || len(enc.Ids) > 512 {
			return ""
		}
		return context
	}

	// Extract emails.
	for _, loc := range emailRegex.FindAllStringIndex(text, -1) {
		start, end := loc[0], loc[1]
		contacts = append(contacts, Contact{
			Value:   text[start:end],
			Context: getContext(start, end),
			Type:    ContactTypeEmail,
		})
	}

	// Extract phone numbers.
	for _, loc := range phoneRegex.FindAllStringIndex(text, -1) {
		start, end := loc[0], loc[1]
		contacts = append(contacts, Contact{
			Value:   text[start:end],
			Context: getContext(start, end),
			Type:    ContactTypePhone,
		})
	}

	// Extract websites.
	for _, loc := range websiteRegex.FindAllStringIndex(text, -1) {
		start, end := loc[0], loc[1]
		value := text[start:end]
		// Remove trailing punctuation (like a period, comma, semicolon, or colon)
		value = strings.TrimRight(value, ".,;:")
		contacts = append(contacts, Contact{
			Value:   value,
			Context: getContext(start, end),
			Type:    ContactTypeWebsite,
		})
	}

	result = contacts
	return result, nil
}
