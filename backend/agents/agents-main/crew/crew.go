package crew

import (
	"context"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/llm"
)

type Crew struct {
	Agents     []agent.Agent
	ManagerLLM llm.LLM
}

func NewCrew(agents []agent.Agent, managerLLM llm.LLM) *Crew {
	if managerLLM != nil {
		panic("managerLLM logic not implemented yet")
	}

	return &Crew{Agents: agents, ManagerLLM: managerLLM}
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
