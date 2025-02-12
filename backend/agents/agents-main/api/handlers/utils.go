package handlers

import (
	"encoding/json"

	"github.com/segp/agents-main/storage"
)

func createStatusFromEvent(event storage.AgentEvent) *ChatCompletionV2StatusResponse {
	switch event.Type {
	case "error":
		return newErrorResponse(event.Metadata.(map[string]string)["error"])

	case "tool_call_choice":
		metadata := event.Metadata.(map[string]map[string]string)
		toolCallArgsStr := metadata["toolCallChoice"]["arguments"]
		toolCallArgs := map[string]string{}
		json.Unmarshal([]byte(toolCallArgsStr), &toolCallArgs)

		return newPendingResponse(toolCallArgs["description_of_action"])

	case "answer_success":
		metadata := event.Metadata.(map[string]string)
		return newCompletedResponse(metadata["answer"])
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
