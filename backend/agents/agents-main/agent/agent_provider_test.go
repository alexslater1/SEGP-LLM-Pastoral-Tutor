package agent

import (
	"fmt"
	"testing"

	"github.com/segp/agents-main/clock"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/imperial_apis"
	abc_api "github.com/segp/agents-main/imperial_apis/abc-api"
	emarking_api "github.com/segp/agents-main/imperial_apis/emarking-api"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

var (
	store               = storage.NewMemoryStorage()
	agentLlm            = llm.NewMockLLM()
	abc_api_client      = abc_api.NewMockAbcApiClient()
	emarking_api_client = emarking_api.NewMockEmarkingApiClient()
	c                   = clock.NewMockClock()
	h                   = history.NewLocalHistory()
	googleSearchClient  = googleSearch.NewMockGoogleSearchClient()
	searchKnowledge     = knowledge.NewLocalKnowledge()
	imperialApiHandler  = imperial_apis.NewImperialApiHandler(abc_api_client, emarking_api_client)
	config              = storage.AgentConfig{
		Name:        "test",
		Description: "test",
		Prompt:      "test",
		Apis: storage.AgentConfigApis{
			AbcApis:      []string{"get-modules"},
			EmarkingApis: []string{"get-exercises"},
		},
		Tools: []string{string(ToolNameSearch)},
	}
)

func TestAgentProvider(t *testing.T) {
	ap := NewAgentProvider(store, agentLlm, c, h, googleSearchClient, searchKnowledge, imperialApiHandler)
	k := ap.apiKnowledgeFrom(config)
	fmt.Printf("%+v\n", k)
}

func TestAgentRefresh(t *testing.T) {
	ap := NewAgentProvider(store, agentLlm, c, h, googleSearchClient, searchKnowledge, imperialApiHandler)
	assert.Equal(t, 0, len(ap.currentAgents))

	storage.Store(store, config)
	ap.RefreshAgents()
	assert.Equal(t, 1, len(ap.currentAgents))
}

func TestToolHandlerFrom(t *testing.T) {
	ap := NewAgentProvider(store, agentLlm, c, h, googleSearchClient, searchKnowledge, imperialApiHandler)
	toolHandler := ap.toolHandlerFrom(config)
	assert.Equal(t, 2, len(toolHandler.Tools)) // as no_tool counts as tool
}

func TestAgentSetAgents(t *testing.T) {
	storage.Store(store, config)
	ap := NewAgentProvider(store, agentLlm, c, h, googleSearchClient, searchKnowledge, imperialApiHandler)

	assert.Equal(t, 1, len(ap.currentAgents))

	agentProviderConfigs := make([]AgentProviderAgentConfig, 2)
	agentProviderConfigs[0] = AgentProviderAgentConfig{
		Name:         "test",
		Description:  "test",
		Prompt:       "test",
		ToolNames:    []string{string(ToolNameSearch)},
		AbcApis:      []string{"get-modules"},
		EmarkingApis: []string{"get-exercises"},
	}

	agentProviderConfigs[1] = AgentProviderAgentConfig{
		Name:         "test2",
		Description:  "test2",
		Prompt:       "test2",
		ToolNames:    []string{string(ToolNameSearch)},
		AbcApis:      []string{"get-modules"},
		EmarkingApis: []string{"get-exercises"},
	}

	ap.SetAgentConfigs(agentProviderConfigs)
	assert.Equal(t, 2, len(ap.currentAgents))

	ap.RefreshAgents()
	assert.Equal(t, 2, len(ap.currentAgents))
}
