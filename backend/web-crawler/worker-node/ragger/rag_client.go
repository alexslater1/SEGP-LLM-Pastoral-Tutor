package ragger

import (
	"log"
	"regexp"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

const (
	contactContextLengthEitherSide = 50
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
	// Compile regular expressions.
	// Email regex: case-insensitive.
	emailRegex, err := regexp.Compile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	if err != nil {
		return nil, err
	}

	// Phone regex: a simple pattern for US-like phone numbers.
	phoneRegex, err := regexp.Compile(`(?i)(\+?\d{1,3}[\s\-]?)?(\(?\d{3}\)?[\s\-]?)?\d{3}[\s\-]?\d{4}`)
	if err != nil {
		return nil, err
	}

	// Address regex: a very basic heuristic pattern.
	// It looks for a number followed by words and a street designator.
	addressRegex, err := regexp.Compile(`(?i)\d+\s+\w+(?:\s+\w+)*\s+(?:Street|St\.|Road|Rd\.|Avenue|Ave\.|Boulevard|Blvd\.|Lane|Ln\.|Drive|Dr\.)`)
	if err != nil {
		return nil, err
	}

	// Prepare slice to hold found contacts.
	var contacts []Contact

	// A helper function to extract context around a match.
	extractContext := func(matchStart, matchEnd int) string {
		// Compute start and end indexes for context.
		start := matchStart - contactContextLengthEitherSide
		if start < 0 {
			start = 0
		}
		end := matchEnd + contactContextLengthEitherSide
		if end > len(text) {
			end = len(text)
		}
		return text[start:end]
	}

	// Define a slice of regexes and their corresponding contact types.
	type regexEntry struct {
		re    *regexp.Regexp
		cType ContactType
	}
	regexes := []regexEntry{
		{emailRegex, ContactTypeEmail},
		{phoneRegex, ContactTypePhone},
		{addressRegex, ContactTypeAddress},
	}

	// Loop over each regex and find all matches.
	for _, entry := range regexes {
		// Find all matches with their indices.
		matches := entry.re.FindAllStringIndex(text, -1)
		for _, loc := range matches {
			// loc[0] is start, loc[1] is end of the match.
			matchVal := text[loc[0]:loc[1]]
			context := extractContext(loc[0], loc[1])
			// Clean up context by replacing newlines with spaces.
			context = strings.ReplaceAll(context, "\n", " ")
			contacts = append(contacts, Contact{
				Value:   matchVal,
				Context: context,
				Type:    entry.cType,
			})
		}
	}

	return contacts, nil
}
