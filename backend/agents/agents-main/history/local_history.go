package history

type LocalHistory struct {
	ChatHistory        []string
	MessagesAndActions []MessagesAndActions
}

func NewLocalHistory() *LocalHistory {
	return &LocalHistory{
		ChatHistory:        []string{},
		MessagesAndActions: []MessagesAndActions{},
	}
}

func (s *LocalHistory) GetMessageHistory(sessionId string) ([]string, error) {
	return s.ChatHistory, nil
}

func (s *LocalHistory) GetMessagesAndActions(sessionId string) ([]MessagesAndActions, error) {
	return s.MessagesAndActions, nil
}

func (s *LocalHistory) AddMessageHistory(messageHistory []string) error {
	s.ChatHistory = append(s.ChatHistory, messageHistory...)
	return nil
}

func (s *LocalHistory) AddMessagesAndActions(messagesAndActions []MessagesAndActions) error {
	s.MessagesAndActions = append(s.MessagesAndActions, messagesAndActions...)
	return nil
}
