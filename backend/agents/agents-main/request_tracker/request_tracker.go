package request_tracker

import (
	"github.com/segp/agents-main/storage"
)

type RequestTracker struct {
	store storage.Storage
}

func NewRequestTracker(store storage.Storage) *RequestTracker {
	return &RequestTracker{store: store}
}

func (rt *RequestTracker) NewRequest(endpoint string) (string, error) {
	request := storage.AgentRequest{
		Endpoint: endpoint,
	}

	data, err := storage.Store(rt.store, request)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) NewCompletionThinkEvent(agentRequestsId string, thoughts *string) (string, error) {
	event := storage.AgentCompletionRequestEvent{
		AgentRequestID: agentRequestsId,
		Type:           storage.AgentCompletionRequestEventTypeThink,
		LLMResponse:    *thoughts,
	}

	data, err := storage.Store(rt.store, event)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) NewCompletionActEvent(agentRequestsId, toolName string, parameters interface{}) (string, error) {
	event := storage.AgentCompletionRequestEvent{
		AgentRequestID: agentRequestsId,
		Type:           storage.AgentCompletionRequestEventTypeAct,
		Metadata: map[string]interface{}{
			"tool_name":  toolName,
			"parameters": parameters,
		},
	}

	data, err := storage.Store(rt.store, event)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) NewCompletionObservationEvent(agentRequestsId string, observations *string) (string, error) {
	event := storage.AgentCompletionRequestEvent{
		AgentRequestID: agentRequestsId,
		Type:           storage.AgentCompletionRequestEventTypeObserve,
		LLMResponse:    *observations,
	}

	data, err := storage.Store(rt.store, event)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) NewCompletionErrorEvent(agentRequestsId string, err error) (string, error) {
	event := storage.AgentCompletionRequestEvent{
		AgentRequestID: agentRequestsId,
		Type:           storage.AgentCompletionRequestEventTypeError,
		Metadata: map[string]interface{}{
			"error": err.Error(),
		},
	}

	data, err := storage.Store(rt.store, event)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) NewCompletionSuccessEvent(agentRequestsId string, response *string, reason *string) (string, error) {
	event := storage.AgentCompletionRequestEvent{
		AgentRequestID: agentRequestsId,
		Type:           storage.AgentCompletionRequestEventTypeError,
		LLMResponse:    *response,
		Metadata: map[string]interface{}{
			"reason": *reason,
		},
	}

	data, err := storage.Store(rt.store, event)
	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (rt *RequestTracker) LatestCompletionObservationEventFor(agentRequestId string) (*storage.AgentCompletionRequestEvent, error) {
	panic("todo")
}
