package history

import (
	"fmt"
	"slices"

	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

type StoreHistory struct {
	storage storage.Storage
}

func NewStoreHistory(storage storage.Storage) *StoreHistory {
	return &StoreHistory{
		storage: storage,
	}
}

func (s *StoreHistory) GetChatHistory(chatId string) ([]string, error) {
	history := []string{}

	agentRequests, err := storage.GetAll[storage.AgentRequest](s.storage, map[string]string{"chat_id": chatId})
	if err != nil {
		return nil, err
	}

	slices.SortFunc(agentRequests, func(a, b storage.AgentRequest) int {
		return a.CreatedAt.Compare(*b.CreatedAt)
	})

	agentEventsListsTasks := utils.DoAsyncList(agentRequests, func(agentRequest storage.AgentRequest) ([]storage.AgentEvent, error) {
		agentEvents, err := storage.GetAll[storage.AgentEvent](s.storage, map[string]string{"request_id": agentRequest.ID})
		if err != nil {
			return nil, err
		}

		return agentEvents, nil
	})

	agentEventsList, err := utils.GetAsyncList(agentEventsListsTasks)
	if err != nil {
		return nil, err
	}

	for i, agentRequest := range agentRequests {
		query := agentRequest.Metadata.(map[string]interface{})["query"].(string)
		answer := s.getChatMessageFromAgentEvents(agentEventsList[i])

		history = append(history, fmt.Sprintf("Query: %s\nResponse: %s", query, answer))
	}

	return history, nil
}

func (s *StoreHistory) getChatMessageFromAgentEvents(agentEvents []storage.AgentEvent) string {
	for _, event := range agentEvents {
		if event.Type == "error" {
			return "[ERROR]"
		}

		if event.Type == "answer_success" {
			return event.Metadata.(map[string]interface{})["answer"].(string)
		}
	}

	return "[PENDING]"
}
