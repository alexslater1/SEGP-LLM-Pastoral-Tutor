package history

type History interface {
	GetChatHistory() ([]string, error)
	AddChatHistory(chatHistory []string) error
}

