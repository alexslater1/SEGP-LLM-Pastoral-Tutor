package crew

import (
	"context"
	"fmt"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/entity"
)

type EntityGraph map[string][]string

type Crew struct {
	EntityGraph EntityGraph
	AgentMap    map[string]agent.Agent
}

func NewCrew(agents []agent.Agent, entityGraph EntityGraph) *Crew {
	agentMap := make(map[string]agent.Agent)
	for _, agent := range agents {
		agentMap[agent.Id()] = agent
	}

	for _, entityIds := range entityGraph {
		for _, entityId := range entityIds {
			if _, ok := agentMap[entityId]; entityId != entity.UserEntityId && !ok {
				panic(fmt.Sprintf("entity %s not found in agent map", entityId))
			}
		}
	}

	return &Crew{EntityGraph: entityGraph, AgentMap: agentMap}
}

func (c *Crew) Run(ctx context.Context, input string) (*string, *string, error) {

	/*
	 1. Pass input into first agent
	 2. Get response from first agent
	 3. If response says to pipe into next agent, pass response into next agent + back to 2 but for new agent

	 4. Return final response
	*/

	panic("not implemented")
}
