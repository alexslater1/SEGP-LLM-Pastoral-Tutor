package ragger

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	ort "github.com/yalue/onnxruntime_go"
)

type TokenizerConfig struct {
	Vocab     map[string]int `json:"vocab"`
	MaxLength int            `json:"max_length"`
	PadToken  int            `json:"pad_token"`
	SepToken  int            `json:"sep_token"`
	ClsToken  int            `json:"cls_token"`
	UnkToken  int            `json:"unk_token"`
}

func loadTokenizer(path string) (*TokenizerConfig, error) {
	data, err := os.ReadFile(path + "/tokenizer_config.json")
	if err != nil {
		return nil, err
	}

	var config TokenizerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default values if not provided in config
	if config.MaxLength == 0 {
		config.MaxLength = 512 // Standard BERT-style max length
	}
	if config.PadToken == 0 {
		config.PadToken = 0 // Common default
	}
	if config.SepToken == 0 {
		config.SepToken = 102 // Common default for BERT-style models
	}
	if config.ClsToken == 0 {
		config.ClsToken = 101 // Common default for BERT-style models
	}
	if config.UnkToken == 0 {
		config.UnkToken = 100 // Common default
	}

	// Validate vocab is not empty
	if len(config.Vocab) == 0 {
		return nil, fmt.Errorf("vocabulary is empty in tokenizer config")
	}

	return &config, nil
}

type DocumentChunker struct {
	tokenizer        *TokenizerConfig
	paragraphSep     string
	chunkSize        int
	separator        string
	secondaryChunkRe *regexp.Regexp
	chunkOverlap     int
}

func NewDocumentChunker(tokenizer *TokenizerConfig) *DocumentChunker {
	return &DocumentChunker{
		tokenizer:        tokenizer,
		paragraphSep:     "\n\n",
		chunkSize:        256,
		separator:        ". ",
		secondaryChunkRe: regexp.MustCompile(`(?:[.!?]|[\n]{2,})`),
		chunkOverlap:     50,
	}
}

func (dc *DocumentChunker) tokenizeText(text string) int {
	// Use the BGE tokenizer to get actual token count
	tokens, _ := dc.tokenizer.Tokenize(text)
	// Return actual token count (including special tokens since they count towards the limit)
	return len(tokens)
}

func (dc *DocumentChunker) ChunkDocument(text string) []string {
	// Split into paragraphs
	paragraphs := strings.Split(text, dc.paragraphSep)
	var allChunks []string
	var currentChunk strings.Builder
	currentTokenCount := 0

	for i, paragraph := range paragraphs {
		// Skip empty paragraphs
		if strings.TrimSpace(paragraph) == "" {
			continue
		}

		// Calculate tokens for this paragraph
		paragraphTokenCount := dc.tokenizeText(paragraph)

		// If adding this paragraph would exceed chunk size, save current chunk and start new one
		if currentTokenCount > 0 && currentTokenCount+paragraphTokenCount > dc.chunkSize {
			if currentChunk.Len() > 0 {
				allChunks = append(allChunks, currentChunk.String())
				currentChunk.Reset()
				currentTokenCount = 0
			}
		}

		// Add paragraph to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(dc.paragraphSep)
		}
		currentChunk.WriteString(paragraph)
		currentTokenCount += paragraphTokenCount

		// If this is the last paragraph or current chunk is getting large, save it
		if i == len(paragraphs)-1 || currentTokenCount >= dc.chunkSize {
			if currentChunk.Len() > 0 {
				allChunks = append(allChunks, currentChunk.String())
				currentChunk.Reset()
				currentTokenCount = 0
			}
		}
	}

	// Add any remaining content
	if currentChunk.Len() > 0 {
		allChunks = append(allChunks, currentChunk.String())
	}

	return allChunks
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Add this new function to embed text
func embedText(session *ort.AdvancedSession, text string, tokenizer *TokenizerConfig, modelPath string) ([]float32, error) {
	// Use our tokenizer properly
	inputIds, attentionMask := tokenizer.Tokenize(text)

	// Convert to int64 for ONNX
	inputIds64 := make([]int64, len(inputIds))
	attentionMask64 := make([]int64, len(attentionMask))
	for i := range inputIds {
		inputIds64[i] = int64(inputIds[i])
		attentionMask64[i] = int64(attentionMask[i])
	}

	inputShape := ort.NewShape(1, int64(len(inputIds))) // [batch_size, sequence_length]

	inputIdsTensor, err := ort.NewTensor(inputShape, inputIds64)
	if err != nil {
		return nil, err
	}
	defer inputIdsTensor.Destroy()

	attentionMaskTensor, err := ort.NewTensor(inputShape, attentionMask64)
	if err != nil {
		return nil, err
	}
	defer attentionMaskTensor.Destroy()

	// BERT output shape will be [batch_size, sequence_length, hidden_size]
	outputShape := ort.NewShape(1, int64(len(inputIds)), 384) // 384 is hidden size for bge-small
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, err
	}
	defer outputTensor.Destroy()

	// Create new session for each run
	session, err = ort.NewAdvancedSession(modelPath,
		[]string{"input_ids", "attention_mask"},
		[]string{"last_hidden_state"},
		[]ort.Value{inputIdsTensor, attentionMaskTensor},
		[]ort.Value{outputTensor},
		nil)
	if err != nil {
		return nil, err
	}
	defer session.Destroy()

	// Run the model
	err = session.Run()
	if err != nil {
		return nil, err
	}

	// Get embeddings from the output tensor
	embeddings := outputTensor.GetData()

	// Return only the CLS token embedding (first token's embedding)
	clsEmbedding := embeddings[:384] // First 384 values represent CLS token embedding
	return clsEmbedding, nil
}

type Contact struct {
	Value   string
	Context string
	Type    string
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

// Add this method to TokenizerConfig
func (t *TokenizerConfig) Tokenize(text string) ([]int, []int) {
	text = strings.TrimSpace(text)
	if text == "" {
		// Return minimal valid tokens for empty text
		return []int{t.ClsToken, t.SepToken}, []int{1, 1}
	}

	tokenIds := []int{t.ClsToken}

	// Split on whitespace and punctuation
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	for _, word := range words {
		word = strings.ToLower(word)
		if tokenId, exists := t.Vocab[word]; exists {
			tokenIds = append(tokenIds, tokenId)
		} else {
			// For unknown words, add a single UNK token instead of per-character
			tokenIds = append(tokenIds, t.UnkToken)
		}
	}

	tokenIds = append(tokenIds, t.SepToken)

	if len(tokenIds) > t.MaxLength {
		tokenIds = tokenIds[:t.MaxLength]
	}

	attentionMask := make([]int, len(tokenIds))
	for i := range attentionMask {
		attentionMask[i] = 1
	}

	return tokenIds, attentionMask
}
