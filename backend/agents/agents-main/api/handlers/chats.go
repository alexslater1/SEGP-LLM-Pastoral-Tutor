package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
)

type MessageRole string

const (
	MessageRoleUser  MessageRole = "user"
	MessageRoleAgent MessageRole = "agent"
)

type QueryAndResponse struct {
	Type      string   `json:"type"`
	Query     string   `json:"query"`
	RequestID string   `json:"request_id"`
	Actions   []string `json:"actions"`

	Answer        string `json:"answer,omitempty"`
	Error         string `json:"error,omitempty"`
	CurrentAction string `json:"current_action,omitempty"`
}

type ChatHistoryResponse struct {
	Messages []QueryAndResponse `json:"query_and_responses"`
}

func SessionIdsForUser(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := context_keys.GetUserID(r.Context())
		if !ok {
			http.Error(w, "user_id not found in context", http.StatusInternalServerError)
			return
		}

		sessions, err := storage.GetAll[storage.Session](store, storage.NewQueryBuilder().Eq("user_id", userID))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get sessions: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(sessions)
	}
}

func SessionFromId(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.PathValue("session_id")
		session, err := storage.Get[storage.Session](store, sessionId)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get session: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(session)
	}
}

func SetDeletedSessionFromId(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.PathValue("session_id")
        session, err := storage.Update[storage.Session](store, sessionId, map[string]any{"deleted": true})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to set session to deleted: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(session)
	}
}

func ChatHistory(history history.History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := r.PathValue("session_id")
		messageAndActions, err := history.GetMessagesAndActions(chatId)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get message history: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(queryAndResponsesFrom(messageAndActions))
	}
}

func queryAndResponsesFrom(messageAndActions []history.MessagesAndActions) []QueryAndResponse {
	queryAndResponses := []QueryAndResponse{}

	for _, messageAndAction := range messageAndActions {
		queryAndResponses = append(queryAndResponses, QueryAndResponse{
			Type:          string(messageAndAction.Type),
			Query:         messageAndAction.Query,
			RequestID:     messageAndAction.RequestID,
			Actions:       messageAndAction.Actions,
			Answer:        messageAndAction.Answer,
			Error:         messageAndAction.Error,
			CurrentAction: messageAndAction.CurrentAction,
		})
	}

	return queryAndResponses
}
