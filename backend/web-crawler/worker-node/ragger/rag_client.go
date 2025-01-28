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

func (c *RagClient) ChunksFor(text string) ([]string, error) {
	return c.chunker.ChunkDocument(text), nil
}

func (c *RagClient) ContactsFor(text string) ([]Contact, error) {
	return extractContactsWithContext(text, 20, 20), nil
}

func (c *RagClient) EmbeddingsFor(text string) ([]float32, error) {
	return embedText(nil, text, c.tokenizer, c.modelPath+"/model.onnx")
}

func setup(libraryPath string) error {
	if err := os.Setenv("LD_LIBRARY_PATH", libraryPath); err != nil {
		return fmt.Errorf("error loading .env file %v", err)
	}

	ort.SetSharedLibraryPath(libraryPath)
	return ort.InitializeEnvironment()
}
