package history

type History interface {
	GetChatHistory(chatId string) ([]string, error)
}
