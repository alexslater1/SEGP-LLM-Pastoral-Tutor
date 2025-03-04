package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/segp/agents-main/context_keys"
	"github.com/segp/agents-main/history"
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
	llm     llm.LLM
	agents  []Agent
	history history.History

	callback func(event AgentEvent)
}

func NewRouter(llm llm.LLM, agents []Agent, history history.History, callback func(event AgentEvent)) *Router {
	r := &Router{
		llm:      llm,
		agents:   agents,
		history:  history,
		callback: callback,
	}

	return r
}

func (r *Router) Run(ctx context.Context, query string) (*AgentResponse, error) {
	ctx = context_keys.SetAgentID(ctx, r.Id())

	slog.Info("Running router", "query", query)
	handleEvent(r.callback, NewQueryEvent(ctx, query))
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

func (r *Router) chatHistory(ctx context.Context) ([]string, error) {
	sessionId, ok := context_keys.GetSessionID(ctx)
	if !ok {
		return nil, nil
	}
	return r.history.GetMessageHistory(sessionId)
}

func (r *Router) queryAndHistoryStrFrom(ctx context.Context, query string) (string, error) {
	chatHistory, err := r.chatHistory(ctx)
	if err != nil {
		return "", err
	}

	if len(chatHistory) == 0 {
		return fmt.Sprintf("The user has typed the following query: %s", query), nil
	}

	return fmt.Sprintf("Up to this point, the chat history is as follows:\n%s\n\nThe user has now typed the following query: %s", strings.Join(chatHistory, "\n"), query), nil
}

func (r *Router) pickAgentForQuery(ctx context.Context, query string) (Agent, string, error) {
	queryAndHistoryStr, err := r.queryAndHistoryStrFrom(ctx, query)
	if err != nil {
		return nil, "", err
	}

	agentSelectionStr := agentSelectionStringFrom(r.agents)
	prompt := fmt.Sprintf("You are a router for a set of agents. The agents are as follows:\n%s\n\n%s\n\nPlease select the most appropriate agent to answer the query. Your response must specify the agent ID and a reason for the selection. You must always select an agent.", agentSelectionStr, queryAndHistoryStr)
	// slog.Info("Prompt", "prompt", prompt)
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
	handleEvent(r.callback, NewRouterSelectionEvent(ctx, routerSelection.AgentID, routerSelection.Reason))
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
