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

type SetConfigsRequest struct {
	Configs []agent.AgentProviderAgentConfig `json:"configs"`
}

func SetConfigs(agentProvider *agent.AgentProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetConfigsRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest)
			return
		}

		for _, config := range req.Configs {
			if config.Name == "" {
				http.Error(w, "name is required", http.StatusBadRequest)
				return
			}

			if config.Description == "" {
				http.Error(w, "description is required", http.StatusBadRequest)
				return
			}

			if config.Prompt == "" {
				http.Error(w, "prompt is required", http.StatusBadRequest)
				return
			}
		}

		if err = agentProvider.SetAgentConfigs(req.Configs); err != nil {
			http.Error(w, fmt.Sprintf("failed to set agent configs: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type SetAutoUpdatePromptsRequst struct {
	Enabled *bool `json:"enabled"`
}

type SetAutoUpdatePromptsResponse struct {
	Enabled bool `json:"enabled"`
}

func SetAutoUpdatePrompts(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SetAutoUpdatePromptsRequst
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to decode request: %v", err), http.StatusBadRequest)
			return
		}

		if req.Enabled == nil {
			http.Error(w, "enabled is required", http.StatusBadRequest)
			return
		}

		data, err := storage.Store(store, storage.NewFeedbackChecksEnabled(*req.Enabled))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to set auto update prompts: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(SetAutoUpdatePromptsResponse{Enabled: data.Enabled})
	}
}

type GetAutoUpdatePromptsResponse struct {
	Enabled bool `json:"enabled"`
}

func GetAutoUpdatePrompts(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := storage.GetAll[storage.FeedbackChecksEnabled](store, storage.NewQueryBuilder().OrderBy("created_at", storage.OrderByDesc).Limit(1))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get auto update prompts: %v", err), http.StatusInternalServerError)
			return
		}

		if len(data) == 0 {
			json.NewEncoder(w).Encode(GetAutoUpdatePromptsResponse{Enabled: false})
			return
		}

		json.NewEncoder(w).Encode(GetAutoUpdatePromptsResponse{Enabled: data[0].Enabled})
	}
}
