package history

import "fmt"

type StatusResponseType string

const (
	StatusResponseTypePending   StatusResponseType = "pending"
	StatusResponseTypeCompleted StatusResponseType = "completed"
	StatusResponseTypeError     StatusResponseType = "error"
)

type MessagesAndActions struct {
	Type      StatusResponseType `json:"type"`
	Query     string             `json:"query"`
	RequestID string             `json:"request_id"`
	Actions   []string           `json:"actions"`
	AgentIDs  []string           `json:"agent_ids"`

	Answer        string `json:"answer,omitempty"`
	Error         string `json:"error,omitempty"`
	CurrentAction string `json:"current_action,omitempty"`
}

type History interface {
	GetMessageHistory(sessionId string) ([]string, error)
	GetMessagesAndActions(sessionId string) ([]MessagesAndActions, error)
}

func MessagesAndActionsToMessageHistory(messagesAndActions []MessagesAndActions) []string {
	messages := []string{}
	for _, messageAndAction := range messagesAndActions {
		message, response := messageAndResponseFrom(messageAndAction)
		messages = append(messages, fmt.Sprintf("User Message: %s\nAgent Response: %s\n", message, response))
	}
	return messages
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
