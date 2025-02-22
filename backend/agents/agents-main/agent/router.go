package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segp/agents-main/llm"
)

const (
	id          = "router"
	description = "A router for a set of agents"
)

type routerSelectionStructuredOutput struct {
	AgentID string `json:"agent_id"`
	Reason  string `json:"reason"`
}

type Router struct {
	llm    llm.LLM
	agents []Agent

	subscribers []chan AgentEvent
}

func NewRouter(llm llm.LLM, agents []Agent) *Router {
	r := &Router{llm: llm, agents: agents, subscribers: []chan AgentEvent{}}

	for _, agent := range agents {
		go r.forwardEventsLoop(agent.Subscribe())
	}

	return r
}

func (r *Router) Run(ctx context.Context, query string) (*AgentResponse, error) {
	agent, _, err := r.pickAgentForQuery(ctx, query)
	if err != nil {
		return nil, err
	}
	return agent.Run(ctx, query)
}

func (r *Router) Id() string {
	return id
}

func (r *Router) Description() string {
	return description
}

func (r *Router) Subscribe() <-chan AgentEvent {
	ch := make(chan AgentEvent, defaultSubscriberBufferSize)
	r.subscribers = append(r.subscribers, ch)
	return ch
}

func (r *Router) Unsubscribe(ch <-chan AgentEvent) {
	for i, subscriber := range r.subscribers {
		if subscriber == ch {
			close(subscriber)
			r.subscribers = append(r.subscribers[:i], r.subscribers[i+1:]...)
			break
		}
	}
}

func (r *Router) forwardEventsLoop(agentCh <-chan AgentEvent) {
	for event := range agentCh {
		for _, subscriber := range r.subscribers {
			subscriber <- event
		}
	}
}

func (r *Router) pickAgentForQuery(ctx context.Context, query string) (Agent, string, error) {
	agentSelectionStr := agentSelectionStringFrom(r.agents)
	prompt := fmt.Sprintf("You are a router for a set of agents. The agents are as follows:\n%s\n\nThe user has asked the following question: %s\n\nPlease select the most appropriate agent to answer the question. Your response must specify the agent ID and a reason for the selection. You must always select an agent.", agentSelectionStr, query)
	response, err := r.llm.StructuredOutputCompletion(ctx, prompt, routerSelectionStructuredOutput{})
	if err != nil {
		return nil, "", err
	}
	routerSelection := routerSelectionStructuredOutput{}
	err = json.Unmarshal([]byte(*response), &routerSelection)
	if err != nil {
		return nil, "", err
	}
	agent, ok := agentFrom(routerSelection.AgentID, r.agents)
	if !ok {
		return nil, "", fmt.Errorf("agent not found")
	}
	return agent, routerSelection.Reason, nil
}

func agentFrom(agentID string, agents []Agent) (Agent, bool) {
	for _, agent := range agents {
		if agent.Id() == agentID {
			return agent, true
		}
	}
	return nil, false
}

func agentSelectionStringFrom(agents []Agent) string {
	agentSelectionStr := ""
	for _, agent := range agents {
		agentSelectionStr += fmt.Sprintf("ID: %s\n Description: %s\n", agent.Id(), agent.Description())
	}
	return agentSelectionStr
}
