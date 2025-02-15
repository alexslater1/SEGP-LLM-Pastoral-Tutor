package crew

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/entity"
)

// TODO: ADD PROPER TESTS FOR THIS

type EntityGraph map[entity.Entity][]entity.Entity

type Crew struct {
	agentMap map[string]agent.Agent
}

func NewCrew(entityGraph EntityGraph) *Crew {
	agentMap := make(map[string]agent.Agent)
	for e, entities := range entityGraph {
		a, ok := e.(agent.Agent)
		if !ok {
			panic(fmt.Sprintf("entity %s is not an agent. Keys must be agents", e.Id()))
		}
		agentMap[a.Id()] = agent.NewCrewAgent(a, entities)
	}

	return &Crew{agentMap: agentMap}
}

func (c *Crew) Run(ctx context.Context, input string, startAgentId string) (*string, *string, error) {
	if startAgentId == entity.UserEntity.Id() {
		return nil, nil, fmt.Errorf("user entity cannot be start agent entity")
	}

	if _, ok := c.agentMap[startAgentId]; !ok {
		return nil, nil, fmt.Errorf("start agent entity %s not found in entity graph", startAgentId)
	}

	entityId := startAgentId
	task := input
	for {
		agent, ok := c.agentMap[entityId]
		if !ok {
			return nil, nil, fmt.Errorf("agent %s not found in entity graph", entityId)
		}
		slog.Info("Crew: running agent", "agent", entityId, "task", task)

		agentResponse, err := agent.Run(ctx, task)
		if err != nil {
			return nil, nil, err
		}

		if agentResponse.OffloadTask == nil {
			return agentResponse.Answer, agentResponse.Reason, nil
		}

		if agentResponse.OffloadTask.Entity == entity.UserEntity {
			return &agentResponse.OffloadTask.Task, nil, nil
		}

		entityId = agentResponse.OffloadTask.Entity.Id()
		task = agentResponse.OffloadTask.Task
	}
}
