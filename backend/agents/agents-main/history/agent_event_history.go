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

	messages := []string{}
	for _, messageAndAction := range messagesAndActions {
		message, response := messageAndResponseFrom(messageAndAction)
		messages = append(messages, fmt.Sprintf("User Message: %s\nAgent Response: %s\n", message, response))
	}

	return messages, nil
}

func (h *AgentEventHistory) GetMessagesAndActions(sessionId string) ([]MessagesAndActions, error) {
	agentRequestSessions, err := storage.GetAll[storage.RequestSession](h.store, map[string]string{"session_id": sessionId})
	if err != nil {
		return nil, err
	}

	sort.Slice(agentRequestSessions, func(i, j int) bool {
		return agentRequestSessions[i].CreatedAt.Before(*agentRequestSessions[j].CreatedAt)
	})

	requestEventsTasks := utils.DoAsyncList(agentRequestSessions, func(agentRequestSession storage.RequestSession) ([]storage.AgentEvent, error) {
		agentEvents, err := storage.GetAll[storage.AgentEvent](h.store, map[string]string{"request_id": agentRequestSession.RequestID})
		if err != nil {
			return nil, err
		}

		sort.Slice(agentEvents, func(i, j int) bool {
			return agentEvents[i].CreatedAt.Before(*agentEvents[j].CreatedAt)
		})

		return agentEvents, nil
	})

	requestEventsLists, err := utils.GetAsyncList(requestEventsTasks)
	if err != nil {
		return nil, err
	}

	messagesAndActions := make([]MessagesAndActions, len(agentRequestSessions))
	for i, agentEvents := range requestEventsLists {
		messagesAndActions[i] = *h.messagesAndActionsFrom(agentRequestSessions[i].RequestID, agentEvents)
	}

	return messagesAndActions, nil
}

func (h *AgentEventHistory) messagesAndActionsFrom(requestId string, agentEvents []storage.AgentEvent) *MessagesAndActions {
	var (
		query  *string
		answer string
		err    string
	)

	actions := []string{}
	mostRecentAction := ""

	for _, event := range agentEvents {
		if event.Type == "tool_call_choice" {
			toolCallChoice := event.Metadata.(map[string]interface{})["toolCallChoice"].(map[string]interface{})

			parsedArgs := map[string]string{}
			err := json.Unmarshal([]byte(toolCallChoice["arguments"].(string)), &parsedArgs)
			if err != nil {
				panic(fmt.Sprintf("failed to unmarshal tool call choice arguments: %v", err))
			}

			a := parsedArgs["description_of_action"]
			mostRecentAction = a
			actions = append(actions, a)
			continue
		}

		if event.Type == "tool_call_result" {
			actions = append(actions, "Thinking")
			continue
		}

		if event.Type == "answer_success" {
			answer = event.Metadata.(map[string]interface{})["answer"].(string)
			continue
		}

		if event.Type == "error" {
			err = event.Metadata.(map[string]interface{})["error"].(string)
			continue
		}

		if event.Type == "query" {
			q := event.Metadata.(map[string]interface{})["query"].(string)
			if query == nil {
				// we only care about the first agent's query
				// is a quick fix, realistically would like to get the query from the request table metadata itself
				query = &q
			}
		}
	}

	responseType := responseStatusFrom(answer, err)
	currentAction := ""
	if responseType == StatusResponseTypePending {
		currentAction = mostRecentAction
	}

	return &MessagesAndActions{
		Type:      responseType,
		Query:     *query,
		RequestID: requestId,
		Actions:   actions,

		Answer:        answer,
		Error:         err,
		CurrentAction: currentAction,
	}
}

func responseStatusFrom(answer string, error string) StatusResponseType {
	if error != "" {
		return StatusResponseTypeError
	}

	if answer != "" {
		return StatusResponseTypeCompleted
	}

	return StatusResponseTypePending
}

func messageAndResponseFrom(messageAndAction MessagesAndActions) (string, string) {
	message := messageAndAction.Query

	if messageAndAction.Type == StatusResponseTypeCompleted {
		return message, messageAndAction.Answer
	}

	if messageAndAction.Type == StatusResponseTypePending {
		return message, "[PENDING]"
	}

	if messageAndAction.Type == StatusResponseTypeError {
		return message, "[ERROR]"
	}

	panic(fmt.Sprintf("unknown status response type: %v", messageAndAction.Type))
}
