package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/segp/agents-main/agent"
)

type ChatCompletionRequest struct {
	Query string `json:"query"`
}

type ChatCompletionResponse struct {
	Response string `json:"response"`
	Reason string `json:"reason"`
}

func ChatCompletion() http.HandlerFunc {
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
		
		agent := agent.NewDefaultReActAgent()

		response, reason, err := agent.Run(req.Query)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := ChatCompletionResponse{
			Response: *response,
			Reason: *reason,
		}

		json.NewEncoder(w).Encode(res)
	}
}