package agent

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segp/agents-main/storage"
)

type EventStoringAgent struct {
	Agent
	storage storage.Storage

	subscribers []chan AgentEvent
}

func NewEventStoringAgent(agent Agent, storage storage.Storage) *EventStoringAgent {
	ch := agent.Subscribe()
	a := &EventStoringAgent{
		Agent:   agent,
		storage: storage,

		subscribers: make([]chan AgentEvent, 0),
	}

	go a.storageLoop(ch)
	return a
}

func (a *EventStoringAgent) Run(ctx context.Context, input string) (*AgentResponse, error) {
	return a.Agent.Run(ctx, input)
}

func (a *EventStoringAgent) Subscribe() <-chan AgentEvent {
	ch := make(chan AgentEvent, defaultSubscriberBufferSize)
	a.subscribers = append(a.subscribers, ch)
	return ch
}

func (a *EventStoringAgent) Unsubscribe(ch <-chan AgentEvent) {
	for i, subscriber := range a.subscribers {
		if subscriber == ch {
			close(subscriber)
			a.subscribers = append(a.subscribers[:i], a.subscribers[i+1:]...)
			break
		}
	}
}

func (a *EventStoringAgent) storageLoop(ch <-chan AgentEvent) {
	for event := range ch {
		a.handleStoreEvent(event)
		for _, subscriber := range a.subscribers {
			subscriber <- event
		}
	}
}

func (a *EventStoringAgent) handleStoreEvent(event AgentEvent) {
	storedEvent := storage.NewAgentEvent(event.RequestID, string(event.Type), event.Data)
	d, err := storage.Store(a.storage, storedEvent)
	if err != nil {
		slog.Error("Error storing event", "error", err)
	}

	fmt.Printf("\n\n[Event Storing Agent] Stored event: %+v\n\n", d)
}
