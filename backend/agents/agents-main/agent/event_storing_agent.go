package agent

import (
	"log/slog"

	"github.com/segp/agents-main/storage"
)

type EventStoringAgent struct {
	Agent
	storage storage.Storage
}

func NewEventStoringAgent(agent Agent, storage storage.Storage) *EventStoringAgent {
	return &EventStoringAgent{
		Agent:   agent,
		storage: storage,
	}
}

func (a *EventStoringAgent) Run(input string, requestId string) (*string, *string, error) {
	ch := a.Subscribe()
	defer a.Unsubscribe(ch)

	go a.storageLoop(ch)

	return a.Agent.Run(input, requestId)
}

func (a *EventStoringAgent) Subscribe() <-chan AgentEvent {
	return a.Agent.Subscribe()
}

func (a *EventStoringAgent) Unsubscribe(ch <-chan AgentEvent) {
	a.Agent.Unsubscribe(ch)
}

func (a *EventStoringAgent) storageLoop(ch <-chan AgentEvent) {
	for event := range ch {
		a.handleStoreEvent(event)
	}
}

func (a *EventStoringAgent) handleStoreEvent(event AgentEvent) {
	storedEvent := storage.NewAgentEvent(event.RequestID, string(event.Type), event.Data)
	_, err := storage.Store(a.storage, storedEvent)
	if err != nil {
		slog.Error("Error storing event", "error", err)
	}
}
