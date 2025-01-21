package knowledge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/segp/agents-main/utils"
)

const (
	defaultNumChunks           = 5
	defaultSimilarityThreshold = 0.5
)

type ChunkResponse struct {
	ChunkID    int     `json:"chunk_id"`
	Text       string  `json:"text"`
	DocID      int     `json:"doc_id"`
	PosInDoc   int     `json:"pos_in_doc"`
	Similarity float64 `json:"similarity"`
}

type RAGResponse struct {
	Response []ChunkResponse `json:"response"`
}

type RAGKnowledge struct {
	URL string
}

func NewRAGKnowledge(url string) *RAGKnowledge {
	return &RAGKnowledge{URL: utils.Required(url, "RAG_BASE_URL")}
}

func (r *RAGKnowledge) Get(query string) (*string, error) {
	chunks, err := r.chunksFrom(query)
	if err != nil {
		return nil, err
	}

	var texts []string
	for _, chunk := range chunks {
		texts = append(texts, chunk.Text)
	}

	jsonTexts, err := json.Marshal(texts)
	if err != nil {
		return nil, err
	}

	result := string(jsonTexts)
	return &result, nil
}

func (r *RAGKnowledge) chunksFrom(query string) ([]ChunkResponse, error) {
	url := fmt.Sprintf("%s/rag?query=%s&num_chunks=%d&similarity_threshold=%f",
		r.URL,
		url.QueryEscape(query),
		defaultNumChunks,
		defaultSimilarityThreshold)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ragResponse RAGResponse
	if err := json.NewDecoder(resp.Body).Decode(&ragResponse); err != nil {
		return nil, err
	}

	return ragResponse.Response, nil
}
