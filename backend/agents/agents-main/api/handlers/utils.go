package handlers

import (
	"encoding/json"

	"github.com/segp/agents-main/storage"
)

func createStatusFromEvent(event storage.AgentEvent) *ChatCompletionV2StatusResponse {
	switch event.Type {
	case "error":
		return newErrorResponse(event.Metadata.(map[string]interface{})["error"].(string))

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
		return newPendingResponse(action)

	case "answer_success":
		metadata := event.Metadata.(map[string]interface{})
		return newCompletedResponse(metadata["answer"].(string))
	}

	return newPendingResponse("Thinking")
}

// Helper functions to create responses
func newPendingResponse(action string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:          ChatCompletionV2StatusResponseTypePending,
		CurrentAction: action,
	}
}

func newErrorResponse(err string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:  ChatCompletionV2StatusResponseTypeError,
		Error: err,
	}
}

func newCompletedResponse(answer string) *ChatCompletionV2StatusResponse {
	return &ChatCompletionV2StatusResponse{
		Type:   ChatCompletionV2StatusResponseTypeCompleted,
		Answer: answer,
	}
}
