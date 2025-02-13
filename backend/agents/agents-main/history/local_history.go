package history

type LocalHistory struct {
	ChatHistory []string
}

func NewLocalHistory() *LocalHistory {
	return &LocalHistory{
		ChatHistory: []string{},
	}
}

func (s *LocalHistory) GetMessageHistory(sessionId string) ([]string, error) {
	return s.ChatHistory, nil
}

func (s *LocalHistory) AddMessageHistory(messageHistory []string) error {
	s.ChatHistory = append(s.ChatHistory, messageHistory...)
	return nil
}
