package agent

import (
	"context"

	"github.com/segp/agents-main/entity"
)

type CrewAgent struct {
	agent                     Agent
	entitiesCanOffloadTasksTo []entity.Entity
}

func NewCrewAgent(agent Agent, entities []entity.Entity) Agent {
	clonedAgent := agent.clone()
	clonedAgent.addCanOffloadToEntity(entities...)

	return &CrewAgent{
		agent:                     clonedAgent,
		entitiesCanOffloadTasksTo: entities,
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
	return NewCrewAgent(a.agent.clone(), a.entitiesCanOffloadTasksTo)
}

func (a *CrewAgent) addCanOffloadToEntity(entities ...entity.Entity) {
	a.agent.addCanOffloadToEntity(entities...)
}

func (a *CrewAgent) canOffloadToEntities() []entity.Entity {
	return a.agent.canOffloadToEntities()
}
