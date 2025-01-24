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
