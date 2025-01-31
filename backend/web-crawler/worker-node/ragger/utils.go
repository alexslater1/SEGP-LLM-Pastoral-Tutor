package ragger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sugarme/tokenizer"
	"github.com/yalue/onnxruntime_go"
)

type DocumentChunker struct {
	paragraphSep       string
	maxChunkSizeTokens int
	separator          string
	secondaryChunkRe   *regexp.Regexp
	chunkOverlap       int
	wordsPerToken      float64 // Average words per token approximation
	tokenizer          *tokenizer.Tokenizer
}

func NewDocumentChunker(tokenizer *tokenizer.Tokenizer) *DocumentChunker {
	return &DocumentChunker{
		paragraphSep:       "\n\n",
		maxChunkSizeTokens: 256,
		separator:          ". ",
		secondaryChunkRe:   regexp.MustCompile(`(?:[.!?]|[\n]{2,})`),
		chunkOverlap:       50,
		wordsPerToken:      0.75, // Most tokenizers average about 0.75 tokens per word
		tokenizer:          tokenizer,
	}
}

func (dc *DocumentChunker) ChunkDocument(text string) []string {
	// Split into paragraphs
	paragraphs := strings.Split(text, dc.paragraphSep)
	var allChunks []string
	var currentChunk strings.Builder

	for _, paragraph := range paragraphs {
		// Skip empty paragraphs
		if strings.TrimSpace(paragraph) == "" {
			continue
		}

		// Try adding paragraph to current chunk
		var proposedChunk string
		if currentChunk.Len() > 0 {
			proposedChunk = currentChunk.String() + dc.paragraphSep + paragraph
		} else {
			proposedChunk = paragraph
		}

		// Check actual token count
		encodeInput := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(proposedChunk))
		encoding, err := dc.tokenizer.Encode(encodeInput, true)
		if err != nil {
			// If encoding fails, treat it as exceeding limit to be safe
			if currentChunk.Len() > 0 {
				allChunks = append(allChunks, currentChunk.String())
				currentChunk.Reset()
			}
			continue
		}

		// If adding this paragraph would exceed token limit, save current chunk and start new one
		if len(encoding.Ids) > dc.maxChunkSizeTokens {
			if currentChunk.Len() > 0 {
				allChunks = append(allChunks, currentChunk.String())
				currentChunk.Reset()
			}
			// Try to add paragraph as a new chunk
			encodeInput = tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(paragraph))
			encoding, err = dc.tokenizer.Encode(encodeInput, true)
			if err == nil && len(encoding.Ids) <= dc.maxChunkSizeTokens {
				currentChunk.WriteString(paragraph)
			} else {
				// If single paragraph is too long, first try splitting by sentences
				sentences := dc.secondaryChunkRe.Split(paragraph, -1)
				currentSentences := make([]string, 0)

				for _, sentence := range sentences {
					// If a single sentence is too long, split it into smaller chunks
					encodeInput = tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(sentence))
					encoding, err = dc.tokenizer.Encode(encodeInput, true)
					if err == nil && len(encoding.Ids) > dc.maxChunkSizeTokens {
						// Split long sentence into smaller chunks by word count
						words := strings.Fields(sentence)
						estimatedChunkSize := int(float64(dc.maxChunkSizeTokens) * dc.wordsPerToken)

						for i := 0; i < len(words); i += estimatedChunkSize {
							end := i + estimatedChunkSize
							if end > len(words) {
								end = len(words)
							}
							chunk := strings.Join(words[i:end], " ")
							if len(currentSentences) > 0 {
								allChunks = append(allChunks, strings.Join(currentSentences, ". "))
								currentSentences = []string{}
							}
							currentSentences = append(currentSentences, chunk)
						}
						continue
					}

					// Normal sentence handling
					if len(currentSentences) > 0 {
						testChunk := strings.Join(append(currentSentences, sentence), ". ")
						encodeInput = tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(testChunk))
						encoding, err = dc.tokenizer.Encode(encodeInput, true)
						if err != nil || len(encoding.Ids) > dc.maxChunkSizeTokens {
							// Save current sentences and start new chunk
							allChunks = append(allChunks, strings.Join(currentSentences, ". "))
							currentSentences = []string{sentence}
						} else {
							currentSentences = append(currentSentences, sentence)
						}
					} else {
						currentSentences = append(currentSentences, sentence)
					}
				}
				if len(currentSentences) > 0 {
					currentChunk.WriteString(strings.Join(currentSentences, ". "))
				}
			}
		} else {
			// Paragraph fits, add it to current chunk
			if currentChunk.Len() > 0 {
				currentChunk.WriteString(dc.paragraphSep)
			}
			currentChunk.WriteString(paragraph)
		}
	}

	// Add any remaining content
	if currentChunk.Len() > 0 {
		allChunks = append(allChunks, currentChunk.String())
	}

	return allChunks
}

func (r *RagClient) inputIdsAndAttentionMasks(texts []string) ([][]int, [][]int, error) {
	var inputIds [][]int
	var attentionMasks [][]int
	maxLength := 0

	// First pass: tokenize and find max length
	for _, text := range texts {
		encodeInput := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(text))
		encoding, err := r.tokenizer.Encode(encodeInput, true)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encode text: %v", err)
		}

		inputIds = append(inputIds, encoding.Ids)
		attentionMasks = append(attentionMasks, encoding.AttentionMask)
		if len(encoding.Ids) > maxLength {
			maxLength = len(encoding.Ids)
		}
	}

	// Pad all sequences to maxLength
	for i := range inputIds {
		if len(inputIds[i]) < maxLength {
			padding := make([]int, maxLength-len(inputIds[i]))
			inputIds[i] = append(inputIds[i], padding...)
			attentionMasks[i] = append(attentionMasks[i], padding...)
		}
	}

	return inputIds, attentionMasks, nil
}

func (r *RagClient) tensorsFromInputIdsAndAttentionMasks(inputIds [][]int, attentionMasks [][]int) (onnxruntime_go.Value, onnxruntime_go.Value, error) {
	batchSize := len(inputIds)
	maxLength := len(inputIds[0])

	inputIdsTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(int64(batchSize), int64(maxLength)),
		make([]int64, batchSize*maxLength),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create input_ids tensor: %v", err)
	}

	attentionMaskTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(int64(batchSize), int64(maxLength)),
		make([]int64, batchSize*maxLength),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create attention_mask tensor: %v", err)
	}

	// Copy values into tensors
	for i := 0; i < batchSize; i++ {
		for j := 0; j < maxLength; j++ {
			idx := i*maxLength + j
			inputIdsTensor.GetData()[idx] = int64(inputIds[i][j])
			attentionMaskTensor.GetData()[idx] = int64(attentionMasks[i][j])
		}
	}

	return inputIdsTensor, attentionMaskTensor, nil
}

func (r *RagClient) embedBatch(texts []string) ([][]float32, error) {
	inputIds, attentionMasks, err := r.inputIdsAndAttentionMasks(texts)
	fmt.Println(len(inputIds[0]))
	if err != nil {
		return nil, fmt.Errorf("failed to get inputIds and attentionMasks: %v", err)
	}

	inputIdsTensor, attentionMaskTensor, err := r.tensorsFromInputIdsAndAttentionMasks(inputIds, attentionMasks)
	if err != nil {
		return nil, fmt.Errorf("failed to get tensors from inputIds and attentionMasks: %v", err)
	}
	defer inputIdsTensor.Destroy()
	defer attentionMaskTensor.Destroy()

	batchSize := len(inputIds)
	maxLength := len(inputIds[0])

	outputTensor, err := onnxruntime_go.NewTensor[float32](
		onnxruntime_go.NewShape(int64(batchSize), int64(maxLength), 384),
		make([]float32, batchSize*maxLength*384),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %v", err)
	}
	defer outputTensor.Destroy()

	// Set up inputs and outputs for dynamic session
	inputs := []onnxruntime_go.Value{inputIdsTensor, attentionMaskTensor}
	outputs := []onnxruntime_go.Value{outputTensor}
	defer outputs[0].Destroy()

	// Run inference
	err = r.onnx_session.Run(inputs, outputs)
	if err != nil {
		return nil, fmt.Errorf("failed to run inference: %v", err)
	}

	// Extract embeddings for each text (using [CLS] token)
	embeddings := make([][]float32, batchSize)
	for i := 0; i < batchSize; i++ {
		embedding := make([]float32, 384)
		start := i * 384 // We only want the [CLS] token embedding
		copy(embedding, outputTensor.GetData()[start:start+384])
		embeddings[i] = embedding
	}

	return embeddings, nil
}

// Keep the original Embed method but make it use EmbedBatch internally
func (r *RagClient) embed(text string) ([]float32, error) {
	embeddings, err := r.embedBatch([]string{text})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

func extractContactsWithContext(text string, wordsBefore, wordsAfter int) []Contact {
	var contacts []Contact

	emailPattern := `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
	// Simplified phone pattern that still preserves the plus sign
	phonePattern := `(?:\+\s*1\s*)?(?:\(\s*[0-9]{3}\s*\)|[0-9]{3})[-.\s]*[0-9]{3}[-.\s]*[0-9]{4}`
	urlPattern := `\b(?:https?://|www\.)[a-zA-Z0-9._%+-]+\.[a-zA-Z]{2,}(?:/[^\s]*)?\b`

	// Helper function to get context
	getContext := func(fullText string, matchStart, matchEnd int, matchText string) string {
		words := strings.Fields(fullText)
		var startIdx, endIdx int

		// Find the word index where the match occurs
		currentLen := 0
		for i, word := range words {
			newLen := currentLen + len(word) + 1 // +1 for space
			if currentLen <= matchStart && matchStart < newLen {
				startIdx = max(0, i-wordsBefore)
				endIdx = min(len(words), i+wordsAfter+1)
				break
			}
			currentLen = newLen
		}

		contextWords := words[startIdx:endIdx]
		return strings.Join(contextWords, " ")
	}

	// Find emails
	emailRe := regexp.MustCompile(emailPattern)
	for _, match := range emailRe.FindAllStringIndex(text, -1) {
		matchText := text[match[0]:match[1]]
		context := getContext(text, match[0], match[1], matchText)
		contacts = append(contacts, Contact{
			Value:   matchText,
			Context: context,
			Type:    "email",
		})
	}

	// Find phone numbers
	phoneRe := regexp.MustCompile(phonePattern)
	for _, match := range phoneRe.FindAllStringIndex(text, -1) {
		matchText := text[match[0]:match[1]]
		context := getContext(text, match[0], match[1], matchText)
		contacts = append(contacts, Contact{
			Value:   matchText,
			Context: context,
			Type:    "phone",
		})
	}

	// Find URLs
	urlRe := regexp.MustCompile(urlPattern)
	for _, match := range urlRe.FindAllStringIndex(text, -1) {
		matchText := text[match[0]:match[1]]
		context := getContext(text, match[0], match[1], matchText)
		contacts = append(contacts, Contact{
			Value:   matchText,
			Context: context,
			Type:    "url",
		})
	}

	return contacts
}
