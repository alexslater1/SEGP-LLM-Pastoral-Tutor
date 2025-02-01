package ragger

import (
	"encoding/json"
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/yalue/onnxruntime_go"
)

const (
	outputEmbeddingSize = 384
	tensorMaxTokens     = 512
)

type Embedder struct {
	tokenizer *tokenizer.Tokenizer
	session   *onnxruntime_go.DynamicAdvancedSession
	modelPath string
}

func NewEmbedder(modelPath string, libraryPath string, tok *tokenizer.Tokenizer) (*Embedder, error) {
	onnxruntime_go.SetSharedLibraryPath(libraryPath)
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize environment: %v", err)
	}

	// Create session with input/output tensors
	session, err := onnxruntime_go.NewDynamicAdvancedSession(
		modelPath,
		[]string{"input_ids", "attention_mask"},
		[]string{"last_hidden_state"},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &Embedder{
		tokenizer: tok,
		session:   session,
		modelPath: modelPath,
	}, nil
}

func (e *Embedder) EmbedAll(texts []string) ([][]float32, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred in EmbedAll. Number of texts: %d\nFirst text: %q\n", len(texts), texts[0])
			if len(texts) > 1 {
				fmt.Printf("Last text: %q\n", texts[len(texts)-1])
			}
			panic(r) // re-panic after printing the context
		}
	}()

	var (
		batchSize      = int64(len(texts))
		inputIds       [][]int64
		attentionMasks [][]int64
	)

	maxTokensEncodingLength := 0
	encodings := make([]*tokenizer.Encoding, len(texts))

	for i, text := range texts {
		encoding, err := e.tokenizer.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(text)), true)
		if err != nil {
			return nil, fmt.Errorf("failed to encode text: %v", err)
		}
		encodings[i] = encoding
		if len(encoding.Ids) > maxTokensEncodingLength {
			maxTokensEncodingLength = len(encoding.Ids)
		}
	}

	if maxTokensEncodingLength > tensorMaxTokens {
		return nil, fmt.Errorf("max tokens encoding length %d is greater than tensor max tokens %d", maxTokensEncodingLength, tensorMaxTokens)
	}

	for _, encoding := range encodings {
		ids := make([]int64, maxTokensEncodingLength)
		copy(ids, toInt64(encoding.Ids))

		mask := make([]int64, maxTokensEncodingLength)
		copy(mask, toInt64(encoding.AttentionMask))

		inputIds = append(inputIds, ids)
		attentionMasks = append(attentionMasks, mask)
	}

	inputIdsTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(batchSize, int64(maxTokensEncodingLength)),
		make([]int64, batchSize*int64(maxTokensEncodingLength)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create input_ids tensor: %v", err)
	}

	attentionMaskTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(batchSize, int64(maxTokensEncodingLength)),
		make([]int64, batchSize*int64(maxTokensEncodingLength)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attention_mask tensor: %v", err)
	}

	for i, ids := range inputIds {
		for j, v := range ids {
			inputIdsTensor.GetData()[i*maxTokensEncodingLength+j] = v
		}
	}

	for i, mask := range attentionMasks {
		for j, v := range mask {
			attentionMaskTensor.GetData()[i*maxTokensEncodingLength+j] = v
		}
	}

	outputTensor, err := onnxruntime_go.NewTensor[float32](
		onnxruntime_go.NewShape(batchSize, int64(maxTokensEncodingLength), outputEmbeddingSize),
		make([]float32, batchSize*int64(maxTokensEncodingLength)*outputEmbeddingSize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %v", err)
	}

	err = e.session.Run([]onnxruntime_go.Value{inputIdsTensor, attentionMaskTensor},
		[]onnxruntime_go.Value{outputTensor})
	if err != nil {
		return nil, fmt.Errorf("failed to run inference: %v", err)
	}

	// Extract only the [CLS] embeddings for each text
	clsEmbeddings := make([]float32, batchSize*outputEmbeddingSize)
	for i := int64(0); i < batchSize; i++ {
		copy(clsEmbeddings[i*outputEmbeddingSize:(i+1)*outputEmbeddingSize],
			outputTensor.GetData()[i*int64(maxTokensEncodingLength)*outputEmbeddingSize:i*int64(maxTokensEncodingLength)*outputEmbeddingSize+int64(outputEmbeddingSize)])
	}

	toReturn := make([][]float32, batchSize)
	for i := range toReturn {
		toReturn[i] = clsEmbeddings[i*outputEmbeddingSize : (i+1)*outputEmbeddingSize]
	}

	return toReturn, nil
}

func (e *Embedder) Embed(text string) ([]float32, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred in Embed. Input text: %q\n", text)
			panic(r)
		}
	}()

	// Tokenize input
	encodeInput := tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(text))
	encoding, err := e.tokenizer.Encode(encodeInput, true)
	if err != nil {
		return nil, fmt.Errorf("failed to encode text: %v", err)
	}
	inputIds := encoding.Ids
	attentionMask := encoding.AttentionMask

	// Create input tensors with int64
	inputIdsTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(1, int64(len(inputIds))),
		make([]int64, len(inputIds)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create input_ids tensor: %v", err)
	}
	// Copy values into tensor
	for i, v := range inputIds {
		inputIdsTensor.GetData()[i] = int64(v)
	}

	attentionMaskTensor, err := onnxruntime_go.NewTensor[int64](
		onnxruntime_go.NewShape(1, int64(len(attentionMask))),
		make([]int64, len(attentionMask)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attention_mask tensor: %v", err)
	}
	// Copy values into tensor
	for i, v := range attentionMask {
		attentionMaskTensor.GetData()[i] = int64(v)
	}

	// Create output tensor with 3D shape [batch_size, sequence_length, hidden_size]
	outputTensor, err := onnxruntime_go.NewTensor[float32](
		onnxruntime_go.NewShape(1, int64(len(inputIds)), outputEmbeddingSize),
		make([]float32, 1*len(inputIds)*outputEmbeddingSize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %v", err)
	}

	err = e.session.Run([]onnxruntime_go.Value{inputIdsTensor, attentionMaskTensor},
		[]onnxruntime_go.Value{outputTensor})
	if err != nil {
		return nil, fmt.Errorf("failed to run inference: %v", err)
	}

	// Get the embedding (first token [CLS] of last hidden state)
	// Since the output is 3D [batch_size, sequence_length, hidden_size],
	// we take the first sequence token's embedding
	embedding := make([]float32, outputEmbeddingSize)
	copy(embedding, outputTensor.GetData()[:outputEmbeddingSize])
	return embedding, nil
}

func printEmbedding(embedding []float32) error {
	json, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("error marshalling embedding: %v", err)
	}

	fmt.Printf("%+v", string(json))
	return nil
}

func toInt64(values []int) []int64 {
	int64Values := make([]int64, len(values))
	for i, v := range values {
		int64Values[i] = int64(v)
	}
	return int64Values
}
