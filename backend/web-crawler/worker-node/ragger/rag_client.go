package ragger

import (
	"fmt"
	"log"

	"github.com/ethanhosier/worker-node/utils"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
	"github.com/yalue/onnxruntime_go"
)

type Contact struct {
	Value   string
	Context string
	Type    string
}

type RagClient struct {
	chunker      *DocumentChunker
	modelPath    string
	onnx_session *onnxruntime_go.DynamicAdvancedSession
	tokenizer    *tokenizer.Tokenizer
}

func NewRagClient(modelPath string, libPath string) *RagClient {
	onnxruntime_go.SetSharedLibraryPath(utils.Required(libPath, "libonx library path"))
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize environment: %v", err)
	}

	session, err := onnxruntime_go.NewDynamicAdvancedSession(modelPath+"/model.onnx",
		[]string{"input_ids", "attention_mask"},
		[]string{"last_hidden_state"},
		nil)
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	tok, err := pretrained.FromFile(modelPath + "/tokenizer.json")
	if err != nil {
		log.Fatalf("failed to load tokenizer: %v", err)
	}

	return &RagClient{
		chunker:      NewDocumentChunker(tok),
		modelPath:    modelPath,
		onnx_session: session,
		tokenizer:    tok,
	}
}

func (c *RagClient) ChunksFrom(text string) ([]string, error) {
	return c.chunker.ChunkDocument(text), nil
}

func (c *RagClient) ContactsFrom(text string) ([]Contact, error) {
	return extractContactsWithContext(text, 20, 20), nil
}

func (c *RagClient) EmbeddingsFor(text string) ([]float32, error) {
	return c.embed(text)
}

func (c *RagClient) EmbeddingsForAll(texts []string) ([][]float32, error) {
	const batchSize = 30 // Process 20 texts at a time to limit memory usage
	var allEmbeddings [][]float32

	for i := 0; i < len(texts); i += batchSize {
		end := i + batchSize
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		batchEmbeddings, err := c.embedBatch(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to embed batch %d-%d: %w", i, end, err)
		}

		allEmbeddings = append(allEmbeddings, batchEmbeddings...)
	}

	return allEmbeddings, nil
}
