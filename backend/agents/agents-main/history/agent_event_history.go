package history

import (
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

	history := make([]string, len(agentRequestSessions))
	for i, requestEvents := range requestEventsLists {
		history[i] = h.getQueryResponseFromAgentEvents(requestEvents)
	}

	return history, nil
}

func (h *AgentEventHistory) getQueryResponseFromAgentEvents(agentEvents []storage.AgentEvent) string {
	var query *string
	var response *string

	for _, event := range agentEvents {
		switch event.Type {
		case "query":
			q := event.Metadata.(map[string]interface{})["query"].(string)
			query = &q
		case "answer_success":
			r := event.Metadata.(map[string]interface{})["answer"].(string)
			response = &r
		case "error":
			r := "[ERROR]"
			response = &r
		}
	}

	if query == nil {
		panic("IMPOSSIBLE STATE: shouldnt get here right?? idk should handle this possibly at some point")
	}

	if response == nil {
		r := "[PENDING]	"
		response = &r
	}

	return fmt.Sprintf("Query: %s\nResponse: %s", *query, *response)
}
