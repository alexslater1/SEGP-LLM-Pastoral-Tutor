package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/segp/agents-main/storage"
)

type MessageV3 struct {
	Role      string                           `json:"role"`
	Content   string                           `json:"content"`
	RequestID string                           `json:"request_id"`
	Statuses  []ChatCompletionV2StatusResponse `json:"statuses"`
}

type ChatHistoryV3Response struct {
	Messages []MessageV3 `json:"messages"`
}

func ChatHistoryV3(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatID := r.PathValue("chat_id")
		if chatID == "" {
			http.Error(w, "chat_id is required", http.StatusBadRequest)
			return
		}

		requests, err := storage.GetAll[storage.AgentRequest](store, map[string]string{"chat_id": chatID})
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting requests %v", err.Error()), http.StatusInternalServerError)
			return
		}

		sort.Slice(requests, func(i, j int) bool {
			return requests[i].CreatedAt.Before(*requests[j].CreatedAt)
		})

		response := []MessageV3{}
		for _, request := range requests {
			request, events, err := getOrderedRequestEvents(store, request.ID)
			if err != nil {
				http.Error(w, fmt.Sprintf("error getting events %v", err.Error()), http.StatusInternalServerError)
				return
			}

			requestMetadata := request.Metadata.(map[string]interface{})
			response = append(response, MessageV3{
				Role:      "user",
				Content:   requestMetadata["query"].(string),
				RequestID: request.ID,
				Statuses:  nil,
			})

			var statuses []ChatCompletionV2StatusResponse
			for _, event := range events {
				statuses = append(statuses, *createStatusFromEvent(event))
			}

			content := "Thinking"
			if len(events) > 0 {
				content = eventContent(events[len(events)-1])
			}

			response = append(response, MessageV3{
				Role:      "assistant",
				Content:   content,
				RequestID: request.ID,
				Statuses:  statuses,
			})
		}

		json.NewEncoder(w).Encode(ChatHistoryV3Response{
			Messages: response,
		})
	}
}

// Helper functions to create responses
func getOrderedRequestEvents(store storage.Storage, requestID string) (*storage.AgentRequest, []storage.AgentEvent, error) {
	events, err := storage.GetAll[storage.AgentEvent](store, map[string]string{"request_id": requestID})
	if err != nil {
		return &storage.AgentRequest{}, nil, err
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(*events[j].CreatedAt)
	})

	request, err := storage.Get[storage.AgentRequest](store, requestID)
	if err != nil {
		return &storage.AgentRequest{}, nil, err
	}

	return request, events, nil
}

func eventContent(event storage.AgentEvent) string {
	switch event.Type {
	case "error":
		return event.Metadata.(map[string]interface{})["error"].(string)

	case "tool_call_choice":
		metadata := event.Metadata.(map[string]interface{})
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
		return action

	case "answer_success":
		metadata := event.Metadata.(map[string]interface{})
		return metadata["answer"].(string)
	}

	return "Thinking"
}
