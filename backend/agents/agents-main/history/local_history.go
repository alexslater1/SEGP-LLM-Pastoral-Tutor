package history

type LocalHistory struct {
	ChatHistory []string
}

func NewLocalHistory() *LocalHistory {
	return &LocalHistory{
		ChatHistory: []string{},
	}
}

func (s *LocalHistory) GetChatHistory(chatId string) ([]string, error) {
	return s.ChatHistory, nil
}

func (s *LocalHistory) AddChatHistory(chatHistory []string) error {
	s.ChatHistory = append(s.ChatHistory, chatHistory...)
	return nil
}
