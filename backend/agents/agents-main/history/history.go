package history

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

	Answer        string `json:"answer,omitempty"`
	Error         string `json:"error,omitempty"`
	CurrentAction string `json:"current_action,omitempty"`
}

type History interface {
	GetMessageHistory(sessionId string) ([]string, error)
	GetMessagesAndActions(sessionId string) ([]MessagesAndActions, error)
}
