package agent

import (
	"context"
	"log"
	"strings"

	"github.com/segp/agents-main/entity"
)

type LoggingAgent struct {
	Agent Agent
}

func NewLoggingAgent(agent Agent) *LoggingAgent {
	return &LoggingAgent{Agent: agent}
}

func (a *LoggingAgent) Run(ctx context.Context, input string) (*AgentResponse, error) {
	ch := a.Subscribe()
	defer a.Unsubscribe(ch)

	go a.loggingLoop(ch)

	return a.Agent.Run(ctx, input)
}

func (a *LoggingAgent) Id() string {
	return a.Agent.Id()
}

func (a *LoggingAgent) Description() string {
	return a.Agent.Description()
}

func (a *LoggingAgent) Subscribe() <-chan AgentEvent {
	return a.Agent.Subscribe()
}

func (a *LoggingAgent) Unsubscribe(ch <-chan AgentEvent) {
	a.Agent.Unsubscribe(ch)
}

func (a *LoggingAgent) clone() Agent {
	return NewLoggingAgent(a.Agent.clone())
}

func (a *LoggingAgent) addCanOffloadToEntity(entities ...entity.Entity) {
	a.Agent.addCanOffloadToEntity(entities...)
}

func (a *LoggingAgent) canOffloadToEntities() []entity.Entity {
	return a.Agent.canOffloadToEntities()
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
