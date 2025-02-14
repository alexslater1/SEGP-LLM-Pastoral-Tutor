package crew

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/entity"
)

type EntityGraph map[string][]string

type Crew struct {
	AgentMap map[string]agent.Agent
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

	return &Crew{AgentMap: crewAgentsMapFrom(entityGraph, agentMap)}
}

func (c *Crew) Run(ctx context.Context, input string, startAgentEntityId string) (*string, *string, error) {

	task := input
	entityId := startAgentEntityId
	for {
		slog.Info("Crew: running agent", "agent", entityId, "task", task)
		agent, ok := c.AgentMap[entityId]
		if !ok {
			return nil, nil, fmt.Errorf("agent %s not found", entityId)
		}

		agentResponse, err := agent.Run(ctx, task)
		if err != nil {
			return nil, nil, err
		}

		if agentResponse.OffloadTask == nil {
			return agentResponse.Answer, agentResponse.Reason, nil
		}

		entityId = agentResponse.OffloadTask.EntityID
		task = agentResponse.OffloadTask.Task
	}
}

func crewAgentsMapFrom(entityGraph EntityGraph, agentMap map[string]agent.Agent) map[string]agent.Agent {
	crewAgentMap := make(map[string]agent.Agent)
	for entityId, agentIds := range entityGraph {
		crewAgentMap[entityId] = agent.NewCrewAgent(agentMap[entityId], agentIds)
	}

	return crewAgentMap
}
