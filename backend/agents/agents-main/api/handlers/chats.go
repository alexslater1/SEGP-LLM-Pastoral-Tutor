package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

type MessageRole string

const (
	MessageRoleUser  MessageRole = "user"
	MessageRoleAgent MessageRole = "agent"
)

type QueryAndResponse struct {
	Type      ChatCompletionV2StatusResponseType `json:"type"`
	Query     string                             `json:"query"`
	RequestID string                             `json:"request_id"`
	Actions   []string                           `json:"actions"`

	Answer        string `json:"answer,omitempty"`
	Error         string `json:"error,omitempty"`
	CurrentAction string `json:"current_action,omitempty"`
}

type ChatHistoryResponse struct {
	Messages []QueryAndResponse `json:"query_and_responses"`
}

func ChatHistory(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := r.PathValue("chat_id")
		agentRequests, err := storage.GetAll[storage.AgentRequest](store, map[string]string{"chat_id": chatId})
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting agent requests %v", err.Error()), http.StatusInternalServerError)
			return
		}

		sort.Slice(agentRequests, func(i, j int) bool {
			return agentRequests[i].CreatedAt.After(*agentRequests[j].CreatedAt)
		})

		eventsTasks := utils.DoAsyncList(agentRequests, func(ar storage.AgentRequest) ([]storage.AgentEvent, error) {
			events, err := storage.GetAll[storage.AgentEvent](store, map[string]string{"request_id": ar.ID})
			if err != nil {
				return nil, err
			}

			sort.Slice(events, func(i, j int) bool {
				return events[i].CreatedAt.After(*events[j].CreatedAt)
			})

			return events, nil
		})

		events, err := utils.GetAsyncList(eventsTasks)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting events %v", err.Error()), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(queryAndResponsesFrom(agentRequests, events))
	}
}

func queryAndResponsesFrom(agentRequests []storage.AgentRequest, events [][]storage.AgentEvent) []QueryAndResponse {
	queryAndResponses := make([]QueryAndResponse, len(agentRequests))
	for i, agentRequest := range agentRequests {
		queryAndResponses[i] = queryAndResponseFrom(agentRequest, events[i])
	}
	return queryAndResponses
}

func queryAndResponseFrom(agentRequest storage.AgentRequest, events []storage.AgentEvent) QueryAndResponse {
	statusResponses := make([]ChatCompletionV2StatusResponse, len(events))
	actions := []string{}
	for i, event := range events {
		statusEvent := createStatusFromEvent(event)
		statusResponses[i] = *statusEvent

		if statusEvent.Type == ChatCompletionV2StatusResponseTypePending {
			actions = append(actions, statusEvent.CurrentAction)
		}
	}

	mostRecentStatus := statusResponses[0]
	query := agentRequest.Metadata.(map[string]interface{})["query"].(string)
	base := QueryAndResponse{
		Type:      mostRecentStatus.Type,
		Query:     query,
		RequestID: agentRequest.ID,
		Actions:   actions,
	}

	switch mostRecentStatus.Type {
	case ChatCompletionV2StatusResponseTypePending:
		base.CurrentAction = mostRecentStatus.CurrentAction

	case ChatCompletionV2StatusResponseTypeCompleted:
		base.Answer = mostRecentStatus.Answer

	case ChatCompletionV2StatusResponseTypeError:
		base.Error = mostRecentStatus.Error
	}

	return base
}
