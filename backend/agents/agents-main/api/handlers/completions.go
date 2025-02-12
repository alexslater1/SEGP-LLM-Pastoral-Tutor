package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sort"

	"github.com/segp/agents-main/agent"

	"github.com/segp/agents-main/storage"
)

type ChatCompletionRequestV1V2 struct {
	Query string `json:"query"`
}

type ChatCompletionRequestV3 struct {
	Query  string `json:"query"`
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
}

type ChatCompletionResponse struct {
	Response string `json:"response"`
	Reason   string `json:"reason"`
}

func ChatCompletion(agent agent.Agent, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequestV1V2
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

type ChatCompletionV2Response struct {
	RequestId string `json:"request_id"`
}

func ChatCompletionV2(agent agent.Agent, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequestV1V2
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

		slog.Info("Created request", "request_id", createdReq.ID)

		go func() {
			response, reasoning, err := agent.Run(req.Query, createdReq.ID)
			if err != nil {
				slog.Error("error running agent", "error", err.Error())
				storage.Store(store, storage.NewAgentEvent(createdReq.ID, "error", map[string]string{"error": err.Error()}))
				return
			}

			slog.Info("agent response", "response", *response, "reason", *reasoning)
		}()

		json.NewEncoder(w).Encode(ChatCompletionV2Response{RequestId: createdReq.ID})
	}
}

type ChatCompletionV2StatusResponseType string

const (
	ChatCompletionV2StatusResponseTypePending   ChatCompletionV2StatusResponseType = "pending"
	ChatCompletionV2StatusResponseTypeCompleted ChatCompletionV2StatusResponseType = "completed"
	ChatCompletionV2StatusResponseTypeError     ChatCompletionV2StatusResponseType = "error"
)

type ChatCompletionV2StatusResponse struct {
	Type          ChatCompletionV2StatusResponseType `json:"type"`
	Error         string                             `json:"error,omitempty"`
	Answer        string                             `json:"answer,omitempty"`
	CurrentAction string                             `json:"current_action,omitempty"`
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

	switch newestEvent.Type {
	case "error":
		return newErrorResponse(newestEvent.Metadata.(map[string]interface{})["error"].(string))

	case "tool_call_choice":
		metadata := newestEvent.Metadata.(map[string]interface{})
		toolCall := metadata["toolCallChoice"].(map[string]interface{})

		// Handle both string and map arguments cases
		var action string
		if argsStr, ok := toolCall["arguments"].(string); ok {
			var args map[string]interface{}
			json.Unmarshal([]byte(argsStr), &args)
			action = args["description_of_action"].(string)
		} else {
			arguments := toolCall["arguments"].(map[string]interface{})
			action = arguments["description_of_action"].(string)
		}
		return newPendingResponse(action)

	case "answer_success":
		metadata := newestEvent.Metadata.(map[string]interface{})
		return newCompletedResponse(metadata["answer"].(string))
	}

	return newPendingResponse("Thinking")
}

func chatCompletionV2StatusResponseFromEvents(events []storage.AgentEvent) *ChatCompletionV2StatusResponse {
	lastEvent := getLatestEvent(events, []string{"error", "tool_call_choice", "answer_success", "observation"})
	if lastEvent == nil {
		return newPendingResponse("Thinking")
	}

	switch lastEvent.Type {
	case "error":
		metadata := lastEvent.Metadata.(map[string]interface{})
		return newErrorResponse(metadata["error"].(string))

	case "observation":
		return newPendingResponse("Thinking")

	case "answer_success":
		metadata := lastEvent.Metadata.(map[string]interface{})
		return newCompletedResponse(metadata["answer"].(string))

	case "tool_call_choice":
		metadata := lastEvent.Metadata.(map[string]interface{})
		toolCall := metadata["toolCallChoice"].(map[string]interface{})

		// Handle both string and map arguments cases
		var action string
		if argsStr, ok := toolCall["arguments"].(string); ok {
			var args map[string]interface{}
			json.Unmarshal([]byte(argsStr), &args)
			action = args["descriptionOfAction"].(string)
		} else {
			arguments := toolCall["arguments"].(map[string]interface{})
			action = arguments["descriptionOfAction"].(string)
		}
		return newPendingResponse(action)
	}

	return newPendingResponse("Thinking")
}

func ChatCompletionV3(agent agent.Agent, store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequestV3
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("error decoding json %v", err.Error()), http.StatusBadRequest)
			return
		}

		if req.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		} else if req.ChatID == "" {
			http.Error(w, "chat ID is required", http.StatusBadRequest)
			return
		} else if req.UserID == "" {
			http.Error(w, "user ID is required", http.StatusBadRequest)
			return
		}

		chat, err := createNewChat(store, req.ChatID, req.UserID)

		if err != nil {
			http.Error(w, fmt.Sprintf("error reading/creating chat %v", err.Error()), http.StatusInternalServerError)
			return
		}
		createdReq, err := storage.Store(store, storage.NewAgentRequest("/completion", map[string]string{"query": req.Query}, chat.ID))

		if err != nil {
			http.Error(w, fmt.Sprintf("error storing request %v", err.Error()), http.StatusInternalServerError)
			return
		}

		slog.Info("Created request", "request_id", createdReq.ID)

		go func() {
			response, reasoning, err := agent.Run(req.Query, createdReq.ID)
			if err != nil {
				slog.Error("error running agent", "error", err.Error())
				storage.Store(store, storage.NewAgentEvent(createdReq.ID, "error", map[string]string{"error": err.Error()}))
				return
			}

			slog.Info("agent response", "response", *response, "reason", *reasoning)
		}()

		json.NewEncoder(w).Encode(ChatCompletionV2Response{RequestId: createdReq.ID})
	}
}

// Helper functions to create responses
func newPendingResponse(action string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:          ChatCompletionV2StatusResponseTypePending,
		CurrentAction: action,
	}
}

func newErrorResponse(err string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:  ChatCompletionV2StatusResponseTypeError,
		Error: err,
	}
}

func newCompletedResponse(answer string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:   ChatCompletionV2StatusResponseTypeCompleted,
		Answer: answer,
	}
}

func getLatestEvent(events []storage.AgentEvent, relevantEventTypes []string) *storage.AgentEvent {
	sortedRelevantEvents := []storage.AgentEvent{}

	for _, event := range events {
		if slices.Contains(relevantEventTypes, event.Type) {
			sortedRelevantEvents = append(sortedRelevantEvents, event)
		}
	}

	if len(sortedRelevantEvents) == 0 {
		return nil
	}

	sort.Slice(sortedRelevantEvents, func(i, j int) bool {
		return sortedRelevantEvents[i].CreatedAt.After(*sortedRelevantEvents[j].CreatedAt)
	})

	return &sortedRelevantEvents[0]
}

func createNewChat(store storage.Storage, chatID string, userID string) (*storage.Chat, error) {
	existingChat, err := storage.Get[storage.Chat](store, chatID)
	if err != nil && err.Error() == "no result returned from supabase" {
		createdChat, err := storage.Store(store, storage.NewChat(chatID, userID, "Temp Title", "private"))
		if err != nil {
			return nil, err
		}
		return createdChat, nil
	} else if err != nil {
		return nil, err
	}
	return existingChat, nil
}
