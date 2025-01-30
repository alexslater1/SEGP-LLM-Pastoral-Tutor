package ragger

import (
	"fmt"
	"log"
	"os"

	"github.com/ethanhosier/worker-node/utils"
	ort "github.com/yalue/onnxruntime_go"
)

type RagClient struct {
	tokenizer *TokenizerConfig
	chunker   *DocumentChunker
	modelPath string
}

func NewRagClient(modelPath string, libPath string) *RagClient {
	if err := setup(utils.Required(libPath, "libonx lihrary path")); err != nil {
		log.Fatalf("error initializing rag client environment %v", err)
	}

	tokenizer, err := loadTokenizer(utils.Required(modelPath, "model path"))
	if err != nil {
		panic(err)
	}

	return &RagClient{
		tokenizer: tokenizer,
		chunker:   NewDocumentChunker(tokenizer),
		modelPath: modelPath,
	}
}

func (c *RagClient) ChunksFrom(text string) ([]string, error) {
	return c.chunker.ChunkDocument(text), nil
}

func (c *RagClient) ContactsFrom(text string) ([]Contact, error) {
	return extractContactsWithContext(text, 20, 20), nil
}

func (c *RagClient) EmbeddingsFor(text string) ([]float32, error) {
	return embedText(text, c.tokenizer, c.modelPath+"/model.onnx")
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
		batchEmbeddings, err := embedMultipleTexts(batch, c.tokenizer, c.modelPath+"/model.onnx")
		if err != nil {
			return nil, fmt.Errorf("failed to embed batch %d-%d: %w", i, end, err)
		}

		allEmbeddings = append(allEmbeddings, batchEmbeddings...)
	}

	return allEmbeddings, nil
}

func setup(libraryPath string) error {
	if err := os.Setenv("LD_LIBRARY_PATH", libraryPath); err != nil {
		return fmt.Errorf("error loading .env file %v", err)
	}

	ort.SetSharedLibraryPath(libraryPath)
	return ort.InitializeEnvironment()
}
