package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/crew"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
)

type ChatCompletionV2StatusResponseType string

const (
	ChatCompletionV2StatusResponseTypePending   ChatCompletionV2StatusResponseType = "pending"
	ChatCompletionV2StatusResponseTypeCompleted ChatCompletionV2StatusResponseType = "completed"
	ChatCompletionV2StatusResponseTypeError     ChatCompletionV2StatusResponseType = "error"
)

type ChatCompletionRequest struct {
	Query     string `json:"query"`
	SessionId string `json:"session_id"`
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
	SessionID string `json:"session_id"`
}

func ChatCompletionV2(crew *crew.Crew, agent agent.Agent, store storage.Storage, history history.History) http.HandlerFunc {
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

		requestId, ok := context_keys.GetRequestID(r.Context())
		if !ok {
			slog.Error("request_id not found in context")
			return
		}

		// TODO: put this in different thread maybe? Idk might break some stuff
		ctx, err := linkSessionToRequest(r.Context(), store)
		if err != nil {
			http.Error(w, fmt.Sprintf("error linking session to request %v", err.Error()), http.StatusInternalServerError)
			return
		}

		go func() {
			answer, reason, err := crew.Run(r.Context(), req.Query, agent.Id())
			if err != nil {
				slog.Error("error running agent", "error", err.Error())
				return
			}

			slog.Info("agent response", "answer", *answer)
			if reason != nil {
				slog.Info("agent response", "reason", *reason)
			}

			rr := storage.NewCompletionResult(requestId, answer, reason, err)
			_, err = storage.Store(store, rr)
			if err != nil {
				slog.Error("error storing request result", "error", err.Error())
			}
		}()

		sessionId, ok := context_keys.GetSessionID(ctx)
		if !ok {
			slog.Error("session_id not found in context")
			return
		}

		json.NewEncoder(w).Encode(ChatCompletionV2Response{RequestId: requestId, SessionID: sessionId})
	}
}

func ChatCompletionV2Status(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.PathValue("request_id")
		if requestId == "" {
			http.Error(w, "request_id is required", http.StatusBadRequest)
			return
		}

		rr, err := storage.GetAll[storage.CompletionResult](store, map[string]string{"request_id": requestId})
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting request result %v", err.Error()), http.StatusInternalServerError)
			return
		}

		if len(rr) == 1 {
			response := handleCompletedRequestResult(rr[0])
			json.NewEncoder(w).Encode(response)
			return
		}

		if len(rr) > 1 {
			http.Error(w, "multiple request results found", http.StatusInternalServerError)
			return
		}

		events, err := storage.GetAll[storage.AgentEvent](store, map[string]string{"request_id": requestId})
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting events %v", err.Error()), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(latestThinkingEventStatusResponse(events))
	}
}

func handleCompletedRequestResult(rr storage.CompletionResult) *ChatCompletionV2StatusResponse {
	if rr.Error != nil {
		return &ChatCompletionV2StatusResponse{
			Type:  ChatCompletionV2StatusResponseTypeError,
			Error: *rr.Error,
		}
	}

	if rr.Result != nil {
		return &ChatCompletionV2StatusResponse{
			Type:   ChatCompletionV2StatusResponseTypeCompleted,
			Answer: *rr.Result,
		}
	}

	panic("impossible state???")
}

func latestThinkingEventStatusResponse(events []storage.AgentEvent) *ChatCompletionV2StatusResponse {
	slices.SortFunc(events, func(a, b storage.AgentEvent) int {
		return b.CreatedAt.Compare(*a.CreatedAt)
	})

	for _, event := range events {
		if event.Type == "tool_call_choice" {
			metadata := event.Metadata.(map[string]interface{})
			toolCallArgsStr := metadata["toolCallChoice"].(map[string]interface{})["arguments"].(string)
			toolCallArgs := map[string]string{}
			json.Unmarshal([]byte(toolCallArgsStr), &toolCallArgs)

			return &ChatCompletionV2StatusResponse{
				Type:          ChatCompletionV2StatusResponseTypePending,
				CurrentAction: toolCallArgs["description_of_action"],
			}
		}
	}

	return &ChatCompletionV2StatusResponse{
		Type:          ChatCompletionV2StatusResponseTypePending,
		CurrentAction: "Thinking",
	}
}

// func chatCompletionV2StatusResponseFromEventsFastAgent(events []storage.AgentEvent) *ChatCompletionV2StatusResponse{
// 	if len(events) == 0 {
// 		return &ChatCompletionV2StatusResponse{
// 			Type: ChatCompletionV2StatusResponseTypePending,
// 		}
// 	}

// 	newestEvent := events[0]
// 	for _, event := range events {
// 		if event.CreatedAt.After(*newestEvent.CreatedAt) {
// 			newestEvent = event
// 		}
// 	}

// 	return createStatusFromEvent(newestEvent)
// }

func linkSessionToRequest(ctx context.Context, store storage.Storage) (context.Context, error) {
	requestID, ok := context_keys.GetRequestID(ctx)
	if !ok {
		return ctx, fmt.Errorf("request_id not found in context")
	}

	userID, ok := context_keys.GetUserID(ctx)
	if !ok {
		return ctx, fmt.Errorf("user_id not found in context")
	}

	sessionID, ok := context_keys.GetSessionID(ctx)
	if !ok {
		session, err := storage.Store(store, storage.NewSession(requestID, userID))
		if err != nil {
			return ctx, fmt.Errorf("error creating session %v", err.Error())
		}

		ctx = context_keys.SetSessionID(ctx, session.ID)
		sessionID = session.ID
	}

	_, err := storage.Store(store, storage.NewRequestSession(sessionID, requestID))
	if err != nil {
		return ctx, fmt.Errorf("error creating request session %v", err.Error())
	}

	return ctx, nil
}
