package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/storage"
)

type ChatCompletionV2StatusResponseType string

const (
	ChatCompletionV2StatusResponseTypePending   ChatCompletionV2StatusResponseType = "pending"
	ChatCompletionV2StatusResponseTypeCompleted ChatCompletionV2StatusResponseType = "completed"
	ChatCompletionV2StatusResponseTypeError     ChatCompletionV2StatusResponseType = "error"
)

type ChatCompletionRequest struct {
	Query  string `json:"query"`
	ChatId string `json:"chat_id"`
	UserId string `json:"user_id"`
}

type ChatCompletionResponse struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
}

type ChatCompletionV2StatusResponse struct {
	Type          ChatCompletionV2StatusResponseType `json:"type"`
	Error         string                             `json:"error,omitempty"`
	Answer        string                             `json:"answer,omitempty"`
	CurrentAction string                             `json:"current_action,omitempty"`
}

type ChatCompletionV2Response struct {
	RequestId string `json:"request_id"`
	ChatId    string `json:"chat_id"`
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

		createdReq, err := storage.Store(store, storage.NewAgentRequest("/completion", map[string]string{"query": req.Query}, ""))
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

func ChatCompletionV2(agent agent.Agent, store storage.Storage) http.HandlerFunc {
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

		chatId := req.ChatId
		if req.ChatId == "" {
			createdChat, err := storage.Store(store, storage.NewChat(req.UserId, nil))
			if err != nil {
				http.Error(w, fmt.Sprintf("error creating chat %v", err.Error()), http.StatusInternalServerError)
				return
			}
			chatId = createdChat.ID
		}

		createdReq, err := storage.Store(store, storage.NewAgentRequest("/completion", map[string]string{"query": req.Query}, chatId))
		if err != nil {
			http.Error(w, fmt.Sprintf("error storing request %v", err.Error()), http.StatusInternalServerError)
			return
		}

		slog.Info("Created request", "request_id", createdReq.ID, "chat_id", chatId)

		go func() {
			response, reasoning, err := agent.Run(req.Query, createdReq.ID)
			if err != nil {
				slog.Error("error running agent", "error", err.Error())
				storage.Store(store, storage.NewAgentEvent(createdReq.ID, "error", map[string]string{"error": err.Error()}))
				return
			}

			slog.Info("agent response", "response", *response, "reason", *reasoning)
		}()

		json.NewEncoder(w).Encode(ChatCompletionV2Response{RequestId: createdReq.ID, ChatId: chatId})
	}
}

func ChatCompletionV2Status(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.PathValue("request_id")
		if requestId == "" {
			http.Error(w, "request_id is required", http.StatusBadRequest)
			return
		}

		events, err := storage.GetAll[storage.AgentEvent](store, map[string]string{"request_id": requestId})
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting events %v", err.Error()), http.StatusInternalServerError)
			return
		}

		response := chatCompletionV2StatusResponseFromEventsFastAgent(events)

		json.NewEncoder(w).Encode(response)
	}
}

func chatCompletionV2StatusResponseFromEventsFastAgent(events []storage.AgentEvent) *ChatCompletionV2StatusResponse {
	if len(events) == 0 {
		return &ChatCompletionV2StatusResponse{
			Type: ChatCompletionV2StatusResponseTypePending,
		}
	}

	newestEvent := events[0]
	for _, event := range events {
		if event.CreatedAt.After(*newestEvent.CreatedAt) {
			newestEvent = event
		}
	}

	return createStatusFromEvent(newestEvent)
}
