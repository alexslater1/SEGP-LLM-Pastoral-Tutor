package crew

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/entity"
)

type EntityGraph map[entity.Entity][]entity.Entity

type Crew struct {
	agentMap map[string]agent.Agent
}

func NewCrew(entityGraph EntityGraph) *Crew {
	agentMap := make(map[string]agent.Agent)
	for e := range entityGraph {
		a, ok := e.(agent.Agent)
		if !ok {
			panic(fmt.Sprintf("entity %s is not an agent. Keys must be agents", e.Id()))
		}
		agentMap[a.Id()] = a
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
		agent := c.agentMap[entityId]
		slog.Info("Crew: running agent", "agent", entityId, "task", task)

		agentResponse, err := agent.Run(ctx, input)
		if err != nil {
			return nil, nil, err
		}

		if agentResponse.OffloadTask == nil {
			return agentResponse.Answer, agentResponse.Reason, nil
		}

		entityId = agentResponse.OffloadTask.Entity.Id()
		task = agentResponse.OffloadTask.Task
	}
}
