package ragger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sugarme/tokenizer"
)

const (
	hardLimitTokenCount = 512
)

const (
	maxTokenChunkSize = hardLimitTokenCount / 2
	// How many tokens of overlap to include between chunks.
	// You might choose a different number depending on your application.
	overlapTokenCount = 50
)

// Chunker holds a reference to the tokenizer so that we can count tokens.
type Chunker struct {
	tokenizer *tokenizer.Tokenizer
}

// NewChunker returns a new Chunker.
func NewChunker(tokenizer *tokenizer.Tokenizer) *Chunker {
	if hardLimitTokenCount < maxTokenChunkSize {
		panic(fmt.Sprintf("maxTokenChunkSize must be less than %v as that is what the model can take", hardLimitTokenCount))
	}

	return &Chunker{
		tokenizer: tokenizer,
	}
}

// Chunk splits the input text into medium sized, overlapping chunks
// while ensuring no chunk ever exceeds maxTokenChunkSize tokens.
func (c *Chunker) Chunk(text string) ([]string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred while processing text: %v\nInput text was: %q\n", r, text)
			panic(r)
		}
	}()

	// Input validation and preprocessing
	if text == "" {
		return []string{" "}, nil
	}

	// Split the text into sentences.
	sentences := splitIntoSentences(text)
	if len(sentences) == 0 {
		return nil, errors.New("no content to chunk")
	}

	var chunks []string
	var currentChunkSentences []string
	currentChunkTokenCount := 0

	// Loop over each sentence
	for _, sentence := range sentences {
		// Count tokens for the sentence.
		enc, err := c.tokenizer.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence)), true)
		if err != nil {
			return nil, err
		}
		sentenceTokenCount := len(enc.Ids)

		// Special-case: if this sentence by itself is longer than maxTokenChunkSize,
		// we split it further using a helper.
		if sentenceTokenCount > maxTokenChunkSize {
			// If we already have some sentences in the current chunk, flush that chunk first.
			if len(currentChunkSentences) > 0 {
				chunks = append(chunks, joinSentences(currentChunkSentences))
				// start new chunk with no overlap.
				currentChunkSentences = nil
				currentChunkTokenCount = 0
			}

			splitSentences, err := c.splitLongSentence(sentence)
			if err != nil {
				return nil, err
			}
			chunks = append(chunks, splitSentences...)
			continue
		}

		// If adding this sentence would exceed our max chunk size, then flush current chunk.
		if currentChunkTokenCount+sentenceTokenCount > maxTokenChunkSize {
			// Flush current chunk.
			chunkText := joinSentences(currentChunkSentences)
			chunks = append(chunks, chunkText)

			// Prepare overlap for the next chunk.
			overlapSentences, newCount := c.getOverlap(currentChunkSentences, overlapTokenCount)
			currentChunkSentences = overlapSentences
			currentChunkTokenCount = newCount
		}

		// Add the current sentence.
		currentChunkSentences = append(currentChunkSentences, sentence)
		currentChunkTokenCount += sentenceTokenCount
	}

	// Flush any remaining sentences.
	if len(currentChunkSentences) > 0 {
		chunks = append(chunks, joinSentences(currentChunkSentences))
	}

	return chunks, nil
}

// splitIntoSentences is a simple (naive) sentence splitter.
func splitIntoSentences(text string) []string {
	// Note: This is very naive. In production you may want to use regex or a NLP library.
	sentences := strings.Split(text, ".")
	var trimmed []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if s != "" {
			// Append the period we lost in the split.
			trimmed = append(trimmed, s+".")
		}
	}
	return trimmed
}

// joinSentences joins sentences with a space.
func joinSentences(sentences []string) string {
	return strings.Join(sentences, " ")
}

// getOverlap selects sentences from the end of the current chunk
// such that the total token count is at least overlapTokenCount (or as close as possible without exceeding the chunk size).
func (c *Chunker) getOverlap(sentences []string, desiredOverlap int) ([]string, int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred in getOverlap. Sentences: `%q`\n", sentences)
			panic(r)
		}
	}()
	var overlapSentences []string
	totalTokens := 0
	// Iterate backwards over the sentences.
	for i := len(sentences) - 1; i >= 0; i-- {
		s := sentences[i]
		enc, err := c.tokenizer.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(s)), true)
		if err != nil {
			// In case of error, skip this sentence.
			continue
		}
		// Prepend the sentence (since we are iterating backwards).
		overlapSentences = append([]string{s}, overlapSentences...)
		totalTokens += len(enc.Ids)
		if totalTokens >= desiredOverlap {
			break
		}
	}
	return overlapSentences, totalTokens
}

// splitLongSentence splits a sentence that exceeds maxTokenChunkSize into smaller pieces.
// For simplicity, we will split at whitespace boundaries.
func (c *Chunker) splitLongSentence(sentence string) ([]string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred in splitLongSentence. Sentence: `%q`\n", sentence)
			panic(r)
		}
	}()
	words := strings.Fields(sentence)
	if len(words) == 0 {
		return nil, errors.New("cannot split an empty sentence")
	}

	var chunks []string
	var currentChunkWords []string
	currentTokenCount := 0

	for _, word := range words {
		// Compute token count for the word.
		enc, err := c.tokenizer.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(word)), true)
		if err != nil {
			return nil, err
		}
		wordTokenCount := len(enc.Ids)

		// If a single word is longer than the max chunk size,
		// then we must split the word itself (this is unusual, but we can simply cut the word).
		if wordTokenCount > maxTokenChunkSize {
			// Flush any current chunk.
			if len(currentChunkWords) > 0 {
				chunks = append(chunks, strings.Join(currentChunkWords, " "))
				currentChunkWords = nil
				currentTokenCount = 0
			}
			// Here we simply cut the word into parts.
			runes := []rune(word)
			for start := 0; start < len(runes); start += maxTokenChunkSize {
				end := start + maxTokenChunkSize
				if end > len(runes) {
					end = len(runes)
				}
				chunks = append(chunks, string(runes[start:end]))
			}
			continue
		}

		// Check if adding the word would exceed the chunk size.
		if currentTokenCount+wordTokenCount > maxTokenChunkSize {
			// Flush current chunk.
			chunks = append(chunks, strings.Join(currentChunkWords, " "))
			// Reset and start new chunk.
			currentChunkWords = []string{word}
			currentTokenCount = wordTokenCount
		} else {
			currentChunkWords = append(currentChunkWords, word)
			currentTokenCount += wordTokenCount
		}
	}

	// Flush any remaining words.
	if len(currentChunkWords) > 0 {
		chunks = append(chunks, strings.Join(currentChunkWords, " "))
	}

	return chunks, nil
}
