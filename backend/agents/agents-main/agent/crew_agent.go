package agent

import "context"

type CrewAgent struct {
	agent                 Agent
	entityIdsCanOffloadTo []string
}

func NewCrewAgent(agent Agent, entityIdsCanOffloadTo []string) Agent {
	clonedAgent := agent.clone()
	clonedAgent.addCanOffloadToEntity(append(agent.canOffloadToEntities(), entityIdsCanOffloadTo...)...)

	return &CrewAgent{
		agent:                 clonedAgent,
		entityIdsCanOffloadTo: clonedAgent.canOffloadToEntities(),
	}
}

func (a *CrewAgent) Run(ctx context.Context, input string) (*AgentResponse, error) {
	return a.agent.Run(ctx, input)
}

func (a *CrewAgent) Id() string {
	return a.agent.Id()
}

func (a *CrewAgent) Description() string {
	return a.agent.Description()
}

func (a *CrewAgent) Subscribe() <-chan AgentEvent {
	return a.agent.Subscribe()
}

func (a *CrewAgent) Unsubscribe(ch <-chan AgentEvent) {
	a.agent.Unsubscribe(ch)
}

func (a *CrewAgent) clone() Agent {
	return NewCrewAgent(a.agent.clone(), a.entityIdsCanOffloadTo)
}

func (a *CrewAgent) addCanOffloadToEntity(entityIds ...string) {
	a.agent.addCanOffloadToEntity(entityIds...)
}

func (a *CrewAgent) canOffloadToEntities() []string {
	return a.entityIdsCanOffloadTo
}
