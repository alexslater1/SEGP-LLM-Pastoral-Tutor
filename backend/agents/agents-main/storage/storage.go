package storage

type Storage interface {
	GetChatHistory() ([]string, error)
	AddChatHistory(chatHistory []string) error
}
