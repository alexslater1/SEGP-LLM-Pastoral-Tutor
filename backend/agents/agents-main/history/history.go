package history

type History interface {
	GetMessageHistory(sessionId string) ([]string, error)
}
