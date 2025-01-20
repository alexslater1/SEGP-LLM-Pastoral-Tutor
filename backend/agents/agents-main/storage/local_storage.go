package storage

type LocalStorage struct {
	ChatHistory []string
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		ChatHistory: []string{},
	}
}

func (s *LocalStorage) GetChatHistory() ([]string, error) {
	return s.ChatHistory, nil
}

func (s *LocalStorage) AddChatHistory(chatHistory []string) error {
	s.ChatHistory = append(s.ChatHistory, chatHistory...)
	return nil
}
