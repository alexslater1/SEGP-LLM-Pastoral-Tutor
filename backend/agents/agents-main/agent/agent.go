package agent

import (
	"context"
	"log"

	"strings"

	"github.com/segp/agents-main/storage"
)

type AgentResponse struct {
	Answer *string
	Reason *string
}

type Agent interface {
	Run(ctx context.Context, input string) (*AgentResponse, error)

	Id() string
	Description() string
}

func handleEvent(callback func(event AgentEvent), event AgentEvent) {
	if callback == nil {
		return
	}

	go callback(event)
}

func newEventStoringLoggingCallback(store storage.Storage) func(event AgentEvent) {
	return func(event AgentEvent) {
		log.Printf("[Event]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
		_, err := storage.Store(store, storage.NewAgentEvent(event.RequestID, string(event.Type), event.Data))
		if err != nil {
			log.Printf("[EventStoringCallback] Error storing event: %v. Continuing...", err)
		}
	}
}

func newLoggingCallback() func(event AgentEvent) {
	return func(event AgentEvent) {
		log.Printf("[Event]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
	}
}
