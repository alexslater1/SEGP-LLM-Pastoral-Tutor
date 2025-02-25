package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/segp/agents-main/agent"
	abc_api "github.com/segp/agents-main/imperial_apis/abc-api"
	emarking_api "github.com/segp/agents-main/imperial_apis/emarking-api"
	"github.com/segp/agents-main/storage"
)

var (
	abc_api_client      = abc_api.NewMockAbcApiClient()
	emarking_api_client = emarking_api.NewMockEmarkingApiClient()
)

func GetAllAgents(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agents, err := storage.GetAll[storage.AgentConfig](store, nil) //TODO: fix race condition of deleting / adding agents

		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get agents: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(agents)
	}
}

type ConfigOptionsResponse struct {
	AbcEndpoints      []abc_api.AbcApiEndpoint           `json:"abc_endpoints"`
	EmarkingEndpoints []emarking_api.EmarkingApiEndpoint `json:"emarking_endpoints"`
	AllToolNames      []agent.ToolName                   `json:"all_tool_names"`
}

func GetConfigOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		json.NewEncoder(w).Encode(ConfigOptionsResponse{
			AbcEndpoints:      abc_api_client.AllApiEndpoints(),
			EmarkingEndpoints: emarking_api_client.AllApiEndpoints(),
			AllToolNames:      agent.AllConfigToolNames(),
		})
	}
}

func SetConfigs(agentProvider *agent.AgentProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req []agent.AgentProviderAgentConfig
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest)
			return
		}

		if err = agentProvider.SetAgentConfigs(req); err != nil {
			http.Error(w, fmt.Sprintf("failed to set agent configs: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
