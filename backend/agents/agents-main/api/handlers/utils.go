package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/segp/agents-main/storage"
)

func createStatusFromEvent(event storage.AgentEvent) *ChatCompletionV2StatusResponse {
	switch event.Type {
	case "error":
		return newErrorResponse(event.Metadata.(map[string]interface{})["error"].(string))

	case "tool_call_choice":
		metadata := event.Metadata.(map[string]interface{})
		toolCallArgsStr := metadata["toolCallChoice"].(map[string]interface{})["arguments"].(string)
		toolCallArgs := map[string]string{}
		json.Unmarshal([]byte(toolCallArgsStr), &toolCallArgs)

		return newPendingResponse(toolCallArgs["description_of_action"])

	case "answer_success":
		fmt.Printf("event.Metadata: %+v\n", event.Metadata)
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
