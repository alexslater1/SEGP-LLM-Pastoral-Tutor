package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/llm"
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

func ChatCompletionV2(agent agent.Agent, store storage.Storage, history history.History, llm llm.LLM) http.HandlerFunc {
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
		newCtx, err := linkSessionToRequest(r.Context(), store, req.Query, llm)
		if err != nil {
			http.Error(w, fmt.Sprintf("error linking session to request %v", err.Error()), http.StatusInternalServerError)
			return
		}

		go func() {
			resp, err := agent.Run(newCtx, req.Query)
			if err != nil {
				slog.Error("error running agent", "error", err.Error())
				rr := storage.NewCompletionResult(requestId, nil, nil, err)
				_, err := storage.Store(store, rr)
				if err != nil {
					slog.Error("error storing request result", "error", err.Error())
				}
				return
			}

			slog.Info("agent response", "answer", *resp.Answer)
			if resp.Reason != nil {
				slog.Info("agent response", "reason", *resp.Reason)
			}

			rr := storage.NewCompletionResult(requestId, resp.Answer, resp.Reason, err)
			d, err := storage.Store(store, rr)
			if err != nil {
				slog.Error("error storing request result", "error", err.Error())
			}

			fmt.Printf("yeahhhh %+v\n", d)
		}()

		sessionId, ok := context_keys.GetSessionID(newCtx)
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

		rr, err := storage.GetAll[storage.CompletionResult](store, storage.NewQueryBuilder().Eq("request_id", requestId))
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

		events, err := storage.GetAll[storage.AgentEvent](store, storage.NewQueryBuilder().Eq("request_id", requestId))
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

// Create new context so that doesnt cancel when request completes
// (So agent can keep running in new thread)
func linkSessionToRequest(oldCtx context.Context, store storage.Storage, query string, llm llm.LLM) (context.Context, error) {
	newContext := context.Background()
	requestID, ok := context_keys.GetRequestID(oldCtx)
	if !ok {
		return oldCtx, fmt.Errorf("request_id not found in context")
	}

	userID, ok := context_keys.GetUserID(oldCtx)
	if !ok {
		return oldCtx, fmt.Errorf("user_id not found in context")
	}

	sessionID, ok := context_keys.GetSessionID(oldCtx)
	if !ok {
		session, err := storage.Store(store, storage.NewSession(requestID, userID))
		if err != nil {
			return oldCtx, fmt.Errorf("error creating session %v", err.Error())
		}

		go handleSetSessionName(llm, store, session.ID, query)
		sessionID = session.ID
	}

	_, err := storage.Store(store, storage.NewRequestSession(sessionID, requestID))
	if err != nil {
		return oldCtx, fmt.Errorf("error creating request session %v", err.Error())
	}

	newContext = context_keys.SetRequestID(newContext, requestID)
	newContext = context_keys.SetUserID(newContext, userID)
	newContext = context_keys.SetSessionID(newContext, sessionID)

	return newContext, nil
}

func handleSetSessionName(llm llm.LLM, store storage.Storage, sessionId string, query string) (string, error) {

	const prompt = `Given the following query "%s", generate a name for the session. The name should be a single sentence that captures the essence of the query. Your response should be just the name and nothing else, with no punctuation at the end`

	resp, err := llm.ChatCompletion(context.TODO(), fmt.Sprintf(prompt, query))
	if err != nil {
		log.Printf("error generating session name %v", err.Error())
		return "", fmt.Errorf("error generating session name %v", err.Error())
	}

	session, err := storage.Update[storage.Session](store, sessionId, map[string]interface{}{"name": *resp})
	if err != nil {
		log.Printf("error updating session name %v", err.Error())
		return "", fmt.Errorf("error updating session name %v", err.Error())
	}

	log.Printf("session name set to %v", *session.Name)
	return *session.Name, nil
}
