package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/segp/agents-main/agent"

	"github.com/segp/agents-main/storage"
)

type ChatCompletionRequest struct {
	Query string `json:"query"`
}

type ChatCompletionResponse struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
}

func ChatCompletion(agent agent.Agent, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("error decoding json %v", err.Error()), http.StatusBadRequest)
			return
		}

		if req.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}

		createdReq, err := storage.Store(store, storage.NewAgentRequest("/completion", map[string]string{"query": req.Query}))
		if err != nil {
			http.Error(w, fmt.Sprintf("error storing request %v", err.Error()), http.StatusInternalServerError)
			return
		}

		response, reasoning, err := agent.Run(req.Query, createdReq.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("error running agent %v", err.Error()), http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(ChatCompletionResponse{Response: *response, Reason: *reasoning})
	}
}
