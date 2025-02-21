package knowledge

import (
	"encoding/json"
	"fmt"
	"log"
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

type ChunksData struct {
	Data []ChunkResponse `json:"data"`
}

type ContactResponse struct {
	ContactID     int     `json:"contact_id"`
	Context       string  `json:"context"`
	DocID         int     `json:"doc_id"`
	Contact       string  `json:"contact"`
	PosInContacts int     `json:"pos_in_contacts"`
	ContactType   string  `json:"contact_type"`
	Similarity    float64 `json:"similarity"`
}

type ContactsData struct {
	Data []ContactResponse `json:"data"`
}

type RAGResponse struct {
	Chunks   ChunksData   `json:"chunks"`
	Contacts ContactsData `json:"contacts"`
}

type RAGKnowledge struct {
	URL string
}

func NewRAGKnowledge(url string) *RAGKnowledge {
	return &RAGKnowledge{URL: utils.Required(url, "RAG_BASE_URL")}
}

func (r *RAGKnowledge) Get(query string) (*string, error) {
	chunks, contacts, err := r.chunksFrom(query)
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

	var contactTexts []string
	for _, contact := range contacts {
		contactTexts = append(contactTexts, contact.Context)
	}

	jsonContactTexts, err := json.Marshal(contactTexts)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("Context Chunks: %s\nContact Chunks: %s", string(jsonTexts), string(jsonContactTexts))
	log.Printf("!!!!!!!!! RAGKnowledge result: %s", result)
	return &result, nil
}

func (r *RAGKnowledge) chunksFrom(query string) ([]ChunkResponse, []ContactResponse, error) {
	url := fmt.Sprintf("%s/rag?query=%s&num_chunks=%d&similarity_threshold=%f",
		r.URL,
		url.QueryEscape(query),
		defaultNumChunks,
		defaultSimilarityThreshold)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	var ragResponse RAGResponse
	if err := json.NewDecoder(resp.Body).Decode(&ragResponse); err != nil {
		return nil, nil, err
	}

	return ragResponse.Chunks.Data, ragResponse.Contacts.Data, nil
}

func (r *RAGKnowledge) Name() string {
	return "rag"
}
