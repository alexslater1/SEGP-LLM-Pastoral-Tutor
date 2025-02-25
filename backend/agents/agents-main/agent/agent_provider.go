package agent

import (
	"log"
	"sync"

	"github.com/segp/agents-main/clock"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/imperial_apis"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

type ToolName string

const (
	ToolNameSearch ToolName = "search"
)

type AgentProviderAgentConfig struct {
	AbcApis      []string `json:"abc_apis"`
	EmarkingApis []string `json:"emarking_apis"`
	ToolNames    []string `json:"tool_names"`
	Name         string   `json:"name"`
	Prompt       string   `json:"prompt"`
	Description  string   `json:"description"`
}

type AgentProvider struct {
	storage            storage.Storage
	llm                llm.LLM
	clock              clock.Clock
	history            history.History
	imperialApiHandler *imperial_apis.ImperialApiHandler

	googleSearchClient googleSearch.GoogleSearchClient
	searchKnowledge    knowledge.Knowledge

	currentAgents   []Agent
	currentAgentsMu sync.RWMutex
}

func NewAgentProvider(storage storage.Storage, llm llm.LLM, clock clock.Clock, history history.History, googleSearchClient googleSearch.GoogleSearchClient, searchKnowledge knowledge.Knowledge, imperialApiHandler *imperial_apis.ImperialApiHandler) *AgentProvider {
	ap := &AgentProvider{storage: storage, llm: llm, clock: clock, history: history, googleSearchClient: googleSearchClient, imperialApiHandler: imperialApiHandler}
	ap.RefreshAgents()
	return ap
}

func (ap *AgentProvider) RefreshAgents() error {
	agentConfigs, err := storage.GetAll[storage.AgentConfig](ap.storage, nil)
	if err != nil {
		return err
	}

	newAgents := make([]Agent, len(agentConfigs))
	for i, agentConfig := range agentConfigs {
		newAgents[i], err = ap.newFastAgentFrom(agentConfig)
		if err != nil {
			return err
		}
	}

	ap.currentAgentsMu.Lock()
	defer ap.currentAgentsMu.Unlock()
	ap.currentAgents = newAgents

	return nil
}

func (ap *AgentProvider) GetAgents() []Agent {
	ap.currentAgentsMu.RLock()
	defer ap.currentAgentsMu.RUnlock()
	return ap.currentAgents
}

func (ap *AgentProvider) SetAgentConfigs(configs []AgentProviderAgentConfig) error {
	panic("not implemented")
}

func (ap *AgentProvider) newFastAgentFrom(config storage.AgentConfig) (*FastAgent, error) {
	return newFastAgent(config.Name, config.Description, config.Prompt, ap.toolHandlerFrom(config), ap.llm, ap.apiKnowledgeFrom(config), ap.clock, ap.history), nil
}

func (ap *AgentProvider) toolHandlerFrom(config storage.AgentConfig) *tools.ToolHandler {
	if len(config.Tools) == 0 {
		return nil
	}

	ts := make([]tools.Tool, len(config.Tools))
	for i, toolName := range config.Tools {
		switch toolName {
		case string(ToolNameSearch):
			ts[i] = tools.NewSearchTool(ap.searchKnowledge, tools.NewGoogleSearchFirstResultsPageContentsTool(ap.googleSearchClient, 3))
		default:
			log.Printf("unknown tool name: %s\n", toolName)
		}
	}

	return tools.NewToolHandler(ts)
}

func (ap *AgentProvider) apiKnowledgeFrom(config storage.AgentConfig) *knowledge.ExtraKnowledge {
	extraKnowledgeSources := make([]interface{}, len(config.Apis.AbcApis)+len(config.Apis.EmarkingApis))
	for i, api := range config.Apis.AbcApis {
		extraKnowledgeSources[i] = ap.imperialApiHandler.AbcEndpointsFor(api)[0]()
	}
	for i, api := range config.Apis.EmarkingApis {
		extraKnowledgeSources[i+len(config.Apis.AbcApis)] = ap.imperialApiHandler.EmarkingEndpointsFor(api)[0]()
	}
	return knowledge.NewExtraKnowledge(extraKnowledgeSources...)
}

func AllConfigToolNames() []ToolName {
	return []ToolName{ToolNameSearch}
}
