package history

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

type AgentEventHistory struct {
	store storage.Storage
}

func NewAgentEventHistory(storage storage.Storage) *AgentEventHistory {
	return &AgentEventHistory{
		store: storage,
	}
}

func (h *AgentEventHistory) GetMessageHistory(sessionId string) ([]string, error) {
	messagesAndActions, err := h.GetMessagesAndActions(sessionId)
	if err != nil {
		return nil, err
	}

	return MessagesAndActionsToMessageHistory(messagesAndActions), nil
}

func (h *AgentEventHistory) GetMessagesAndActions(sessionId string) ([]MessagesAndActions, error) {
	agentRequestSessions, err := storage.GetAll[storage.RequestSession](h.store, storage.NewQueryBuilder().Eq("session_id", sessionId))
	if err != nil {
		return nil, err
	}

	sort.Slice(agentRequestSessions, func(i, j int) bool {
		return agentRequestSessions[i].CreatedAt.Before(*agentRequestSessions[j].CreatedAt)
	})

	requestEventsTasks := utils.DoAsyncList(agentRequestSessions, func(agentRequestSession storage.RequestSession) ([]storage.AgentEvent, error) {
		agentEvents, err := storage.GetAll[storage.AgentEvent](h.store, storage.NewQueryBuilder().Eq("request_id", agentRequestSession.RequestID))
		if err != nil {
			return nil, err
		}

		sort.Slice(agentEvents, func(i, j int) bool {
			return agentEvents[i].CreatedAt.Before(*agentEvents[j].CreatedAt)
		})

		return agentEvents, nil
	})

	completionResultsTasks := utils.DoAsyncList(agentRequestSessions, func(agentRequestSession storage.RequestSession) (*storage.CompletionResult, error) {
		resp, err := storage.GetAll[storage.CompletionResult](h.store, storage.NewQueryBuilder().Eq("request_id", agentRequestSession.RequestID))
		if err != nil {
			return nil, err
		}

		if len(resp) == 0 {
			return nil, nil
		}

		return &resp[0], nil
	})

	completionResults, err := utils.GetAsyncList(completionResultsTasks)
	if err != nil {
		return nil, err
	}

	requestEventsLists, err := utils.GetAsyncList(requestEventsTasks)
	if err != nil {
		return nil, err
	}

	messagesAndActions := make([]MessagesAndActions, len(agentRequestSessions))
	for i, agentEvents := range requestEventsLists {
		messagesAndActions[i] = *h.messagesAndActionsFrom(agentRequestSessions[i].RequestID, agentEvents, completionResults[i])
	}

	return messagesAndActions, nil
}

func (h *AgentEventHistory) messagesAndActionsFrom(requestId string, agentEvents []storage.AgentEvent, completionResult *storage.CompletionResult) *MessagesAndActions {
	var (
		query  = ""
		err    string
		answer string
	)

	actions := []string{}
	mostRecentAction := ""

	for _, event := range agentEvents {
		switch event.Type {
		case "tool_call_choice":
			toolCallChoice := event.Metadata.(map[string]interface{})["toolCallChoice"].(map[string]interface{})

			parsedArgs := map[string]string{}
			err := json.Unmarshal([]byte(toolCallChoice["arguments"].(string)), &parsedArgs)
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal tool call choice arguments: %v", err))
			}

			a := parsedArgs["description_of_action"]
			mostRecentAction = a
			actions = append(actions, a)

		case "tool_call_result":
			actions = append(actions, "Thinking")

		// case "answer_success":
		// 	answer = event.Metadata.(map[string]interface{})["answer"].(string)
		// case "answer_error":
		// 	err = event.Metadata.(map[string]interface{})["error"].(string)

		case "query":
			q := event.Metadata.(map[string]interface{})["query"].(string)
			if query == "" {
				// we only care about the first agent's query
				// is a quick fix, realistically would like to get the query from the request table metadata itself
				query = q
			}
		}
	}

	responseType := responseStatusFrom(completionResult, err)
	currentAction := ""
	if responseType == StatusResponseTypePending {
		currentAction = mostRecentAction
	}

	if completionResult != nil {
		if completionResult.Error != nil {
			err = *completionResult.Error
		}

		if completionResult.Result != nil {
			answer = *completionResult.Result
		}
	}

	agentID := agentEvents[0].Metadata.(map[string]interface{})["agentID"].(string)

	return &MessagesAndActions{
		Type:          responseType,
		Query:         query,
		RequestID:     requestId,
		Actions:       actions,
		AgentID:       agentID,
		Answer:        answer,
		Error:         err,
		CurrentAction: currentAction,
	}
}

func responseStatusFrom(completionResult *storage.CompletionResult, error string) StatusResponseType {
	if error != "" {
		return StatusResponseTypeError
	}

	if completionResult != nil && completionResult.Result != nil {
		return StatusResponseTypeCompleted
	}

	return StatusResponseTypePending
}
