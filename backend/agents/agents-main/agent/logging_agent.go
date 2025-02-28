package agent

import (
	"context"
	"log"
	"strings"
)

type LoggingAgent struct {
	Agent       Agent
	subscribers []chan AgentEvent
}

func NewLoggingAgent(agent Agent) *LoggingAgent {
	ch := agent.Subscribe()
	a := &LoggingAgent{
		Agent:       agent,
		subscribers: make([]chan AgentEvent, 0),
	}

	go a.loggingLoop(ch)
	return a
}

func (a *LoggingAgent) Run(ctx context.Context, input string) (*AgentResponse, error) {
	return a.Agent.Run(ctx, input)
}

func (a *LoggingAgent) Id() string {
	return a.Agent.Id()
}

func (a *LoggingAgent) Description() string {
	return a.Agent.Description()
}

func (a *LoggingAgent) Subscribe() <-chan AgentEvent {
	ch := make(chan AgentEvent, defaultSubscriberBufferSize)
	a.subscribers = append(a.subscribers, ch)
	return ch
}

func (a *LoggingAgent) Unsubscribe(ch <-chan AgentEvent) {
	for i, subscriber := range a.subscribers {
		if subscriber == ch {
			close(subscriber)
			a.subscribers = append(a.subscribers[:i], a.subscribers[i+1:]...)
			break
		}
	}
}

func (a *LoggingAgent) close() {
	for _, subscriber := range a.subscribers {
		close(subscriber)
	}
}

func (a *LoggingAgent) loggingLoop(ch <-chan AgentEvent) {
	for event := range ch {
		a.handleLoggingEvent(event)
		for _, subscriber := range a.subscribers {
			subscriber <- event
		}
	}

	a.close()
	a.Agent.Unsubscribe(ch)
}

func (a *LoggingAgent) handleLoggingEvent(event AgentEvent) {
	// if event.Type == AgentEventTypeToolCallResult {
	// 	log.Printf("%s (%s):\n%+v...\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data["toolCallResult"].(string)[:200])
	// 	return
	// }

	log.Printf("[Logging Agent]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
}
