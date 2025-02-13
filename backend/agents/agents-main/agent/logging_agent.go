package agent

import (
	"context"
	"log"
	"strings"
)

type LoggingAgent struct {
	Agent Agent
}

func NewLoggingAgent(agent Agent) *LoggingAgent {
	return &LoggingAgent{Agent: agent}
}

func (a *LoggingAgent) Run(ctx context.Context, input string) (*string, *string, error) {
	ch := a.Subscribe()
	defer a.Unsubscribe(ch)

	go a.loggingLoop(ch)

	return a.Agent.Run(ctx, input)
}

func (a *LoggingAgent) Subscribe() <-chan AgentEvent {
	return a.Agent.Subscribe()
}

func (a *LoggingAgent) Unsubscribe(ch <-chan AgentEvent) {
	a.Agent.Unsubscribe(ch)
}

func (a *LoggingAgent) loggingLoop(ch <-chan AgentEvent) {
	for event := range ch {
		a.handleLoggingEvent(event)
	}
}

func (a *LoggingAgent) handleLoggingEvent(event AgentEvent) {
	// if event.Type == AgentEventTypeToolCallResult {
	// 	log.Printf("%s (%s):\n%+v...\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data["toolCallResult"].(string)[:200])
	// 	return
	// }

	log.Printf("%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
}
