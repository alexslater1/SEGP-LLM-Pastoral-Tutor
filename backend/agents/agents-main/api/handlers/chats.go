package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/segp/agents-main/history"
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

func ChatHistory(history history.History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := r.PathValue("chat_id")
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
