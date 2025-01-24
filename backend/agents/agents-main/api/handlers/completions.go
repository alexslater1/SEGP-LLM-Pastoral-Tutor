package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/request_tracker"
)

type ChatCompletionRequest struct {
	Query string `json:"query"`
}

type ChatCompletionResponse struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
}

func ChatCompletion(agent agent.Agent) http.HandlerFunc {
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

		response, reason, err := agent.Run(req.Query)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := ChatCompletionResponse{
			Response: *response,
			Reason:   *reason,
		}

		json.NewEncoder(w).Encode(res)
	}
}

func ChatCompletionV2(agent agent.Agent, requestTracker *request_tracker.RequestTracker) http.HandlerFunc {
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

		requestId, err := requestTracker.NewRequest("chatCompletionV2")
		if err != nil {
			http.Error(w, fmt.Sprintf("error creating request %v", err.Error()), http.StatusInternalServerError)
		}

		go func() {
			completion, reason, err := agent.Run(req.Query)
			if err != nil {
				log.Printf("Error while processing request %s:%s", requestId, err.Error())
				requestTracker.NewCompletionErrorEvent(requestId, err)
				return
			}

			_, err = requestTracker.NewCompletionSuccessEvent(requestId, completion, reason)
			if err != nil {
				log.Printf("Error when setting success event for request %s:%s", requestId, err.Error())
			}
		}()
	}
}

type NextStepAction string

const (
	NextStepActionReAct NextStepAction = "reAct"
	NextStepActionReply NextStepAction = "reply"
)

type NextStep struct {
	Action        NextStepAction
	OptionalReply string
}

func handleNextStep(query string) (*NextStep, error) {
	handleNextStepPrompt := `
	You are a helpful assistant. You are given a query and some context and you need to determine the next step to take.

	If you can answer the query confidently from your own knowledge and the given context, and there are no further actions which should be taken, you must reply with the Action being reply and the OptionalReply being your reply to the user

	If you cannot answer the query confidently from your own knowledge and the given context, you must reply with the Action being reAct and the OptionalReply being an empty string

	The reAct agent has access to a plethora of tools to access and deal with real time data, such as Google Search and interactions with other systems.

	Here is the query: %s
	Here is the context: %s
	`

	ragKnowledge := knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
	ragResponseStr, err := ragKnowledge.Get(query)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(handleNextStepPrompt, query, *ragResponseStr)

	llm := llm.NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))
	response, err := llm.StructuredOutputCompletion(context.Background(), prompt, NextStep{})
	if err != nil {
		return nil, err
	}

	var parsedResponse NextStep
	if err := json.Unmarshal([]byte(*response), &parsedResponse); err != nil {
		return nil, err
	}

	return &parsedResponse, nil
}
