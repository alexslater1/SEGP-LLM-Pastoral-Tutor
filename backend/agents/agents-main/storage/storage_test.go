package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Store(t *testing.T) {
	storage := NewMemoryStorage()

	req := AgentRequest{ID: "id1", Endpoint: "endpoint1", UserID: "user1"}
	data, err := Store(storage, req)
	assert.NoError(t, err)
	assert.Equal(t, req.ID, data.ID)
	assert.Equal(t, req.Endpoint, data.Endpoint)
}

func TestStorage_Get(t *testing.T) {
	storage := NewMemoryStorage()

	req := AgentRequest{ID: "id1", Endpoint: "endpoint1"}
	storage.store(req.TableName(), req)

	res, err := Get[AgentRequest](storage, "id1")
	assert.NoError(t, err)
	assert.Equal(t, "endpoint1", res.Endpoint)
	assert.Equal(t, "id1", res.ID)
}

func TestStorage_GetNoId(t *testing.T) {
	storage := NewMemoryStorage()

	req := AgentRequest{Endpoint: "endpoint1"}
	data, err := Store(storage, req)
	assert.NoError(t, err)

	res, err := Get[AgentRequest](storage, data.ID)
	assert.NoError(t, err)
	assert.Equal(t, "endpoint1", res.Endpoint)
	assert.Equal(t, data.ID, res.ID)
}

func TestStorage_StoreAll(t *testing.T) {
	storage := NewMemoryStorage()

	reqs := []AgentRequest{
		{ID: "id1", Endpoint: "endpoint1"},
		{ID: "id2", Endpoint: "endpoint2"},
		{ID: "id3", Endpoint: "endpoint1"},
	}

	data, err := StoreAll(storage, reqs...)
	assert.NoError(t, err)

	for i, req := range data {
		assert.Equal(t, reqs[i].ID, req.ID)
		assert.Equal(t, reqs[i].Endpoint, req.Endpoint)
	}
}

func TestStorage_GetAll(t *testing.T) {
	storage := NewMemoryStorage()

	reqs := []AgentRequest{
		{ID: "id1", Endpoint: "endpoint1"},
		{ID: "id2", Endpoint: "endpoint2"},
		{ID: "id3", Endpoint: "endpoint1"},
	}

	_, err := StoreAll(storage, reqs...)
	assert.NoError(t, err)

	res, err := GetAll[AgentRequest](storage, NewQueryBuilder().Eq("endpoint", "endpoint1"))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))

	assert.Equal(t, reqs[0].ID, res[0].ID)
	assert.Equal(t, reqs[2].ID, res[1].ID)
}

func TestStorage_NewRequestResult(t *testing.T) {
	var (
		result = "result1"
		reason = "reason1"
		err    = errors.New("error1")

		store = NewMemoryStorage()
	)

	rr := NewCompletionResult("request1", &result, &reason, err)
	assert.Equal(t, "request1", rr.RequestID)
	assert.Equal(t, "result1", *rr.Result)
	assert.Equal(t, "reason1", *rr.Reason)
	assert.Equal(t, "error1", *rr.Error)

	res, err := Store(store, rr)
	assert.NoError(t, err)
	assert.Equal(t, rr.RequestID, res.RequestID)
	assert.Equal(t, rr.Result, res.Result)
	assert.Equal(t, rr.Reason, res.Reason)
	assert.Equal(t, rr.Error, res.Error)

	// res, err = Get[CompletionResult](store, res.ID)
	// assert.NoError(t, err)
	// assert.Equal(t, rr.RequestID, res.RequestID)
	// assert.Equal(t, rr.Result, res.Result)
	// assert.Equal(t, rr.Reason, res.Reason)
	// assert.Equal(t, rr.Error, res.Error)
}
