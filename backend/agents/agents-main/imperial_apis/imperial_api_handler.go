package imperial_apis

import (
	abc_api "github.com/segp/agents-main/imperial_apis/abc-api"
	emarking_api "github.com/segp/agents-main/imperial_apis/emarking-api"
)

type ImperialApiHandler struct {
	abcClient      *abc_api.MockAbcApiClient
	emarkingClient *emarking_api.MockEmarkingApiClient
}

func NewImperialApiHandler(abcClient *abc_api.MockAbcApiClient, emarkingClient *emarking_api.MockEmarkingApiClient) *ImperialApiHandler {
	return &ImperialApiHandler{
		abcClient:      abcClient,
		emarkingClient: emarkingClient,
	}
}

func (h *ImperialApiHandler) EmarkingApis() []emarking_api.EmarkingApiEndpoint {
	return h.emarkingClient.AllApiEndpoints()
}

func (h *ImperialApiHandler) AbcApis() []abc_api.AbcApiEndpoint {
	return h.abcClient.AllApiEndpoints()
}

func (h *ImperialApiHandler) EmarkingEndpointsFor(endpoints ...string) []func() any {
	endpointMappings := h.emarkingClient.EndpointMappings()
	results := make([]func() any, 0, len(endpoints))
	for _, endpoint := range endpoints {
		results = append(results, endpointMappings[emarking_api.EmarkingApiEndpoint(endpoint)])
	}
	return results
}

func (h *ImperialApiHandler) AbcEndpointsFor(endpoints ...string) []func() any {
	endpointMappings := h.abcClient.EndpointMappings()
	results := make([]func() any, 0, len(endpoints))
	for _, endpoint := range endpoints {
		results = append(results, endpointMappings[abc_api.AbcApiEndpoint(endpoint)])
	}
	return results
}
