package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseResult[T StorageType](result interface{}) (*T, error) {
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	var res T
	err = json.Unmarshal(jsonResult, &res)
	return &res, err
}

func TestMemoryStorage_Store(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("stores data with existing ID", func(t *testing.T) {
		req := AgentRequest{
			ID:       "test-id",
			Endpoint: "test-endpoint",
		}

		result, err := storage.store(req.TableName(), req)
		assert.NoError(t, err)

		res, err := parseResult[AgentRequest](result)
		assert.NoError(t, err)
		assert.Equal(t, req.ID, res.ID)
		assert.Equal(t, req.Endpoint, res.Endpoint)
	})

	t.Run("generates ID when none provided", func(t *testing.T) {
		req := AgentRequest{
			Endpoint: "test-endpoint",
		}

		result, err := storage.store(req.TableName(), req)
		assert.NoError(t, err)

		res, err := parseResult[AgentRequest](result)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.ID) // ID should be generated
		assert.Equal(t, req.Endpoint, res.Endpoint)
	})

	t.Run("handles nil data", func(t *testing.T) {
		_, err := storage.store(StorageTableNameAgentRequests, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data cannot be nil")
	})

	t.Run("handles numeric ID", func(t *testing.T) {
		type NumericIDStruct struct {
			ID       int    `json:"id"`
			Endpoint string `json:"endpoint"`
		}
		req := NumericIDStruct{
			ID:       123,
			Endpoint: "test-endpoint",
		}

		result, err := storage.store(StorageTableNameAgentRequests, req)
		assert.NoError(t, err)

		res, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, 123, res["id"])
		assert.Equal(t, "test-endpoint", res["endpoint"].(string))
	})

	t.Run("generates numeric ID when original type is numeric", func(t *testing.T) {
		type NumericIDStruct struct {
			ID       int    `json:"id,omitempty"`
			Endpoint string `json:"endpoint"`
		}
		req := NumericIDStruct{
			Endpoint: "test-endpoint",
		}

		result, err := storage.store(StorageTableNameAgentRequests, req)
		assert.NoError(t, err)

		res, ok := result.(map[string]interface{})
		assert.True(t, ok)

		// Check that generated ID is a number
		id, ok := res["id"].(int)
		assert.True(t, ok)
		assert.Greater(t, id, 0)
		assert.LessOrEqual(t, id, 1000000)
		assert.Equal(t, "test-endpoint", res["endpoint"].(string))
	})

	t.Run("preserves string ID type when provided", func(t *testing.T) {
		type StringIDStruct struct {
			ID       string `json:"id"`
			Endpoint string `json:"endpoint"`
		}
		req := StringIDStruct{
			ID:       "test-id",
			Endpoint: "test-endpoint",
		}

		result, err := storage.store(StorageTableNameAgentRequests, req)
		assert.NoError(t, err)

		res, ok := result.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "test-id", res["id"])
		assert.Equal(t, "test-endpoint", res["endpoint"])
	})
}

func TestMemoryStorage_Get(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("retrieves stored data", func(t *testing.T) {
		req := AgentRequest{
			ID:       "test-id",
			Endpoint: "test-endpoint",
		}

		// Store the data first
		_, err := storage.store(req.TableName(), req)
		assert.NoError(t, err)

		// Retrieve the data
		result, err := storage.get(req.TableName(), req.ID)
		assert.NoError(t, err)

		res, err := parseResult[AgentRequest](result)
		assert.NoError(t, err)

		assert.Equal(t, req.ID, res.ID)
		assert.Equal(t, req.Endpoint, res.Endpoint)
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := storage.get(StorageTableNameAgentRequests, "non-existent-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "item not found")
	})
}

func TestMemoryStorage_StoreAll(t *testing.T) {
	storage := NewMemoryStorage()

	t.Run("stores multiple items", func(t *testing.T) {
		reqs := []interface{}{
			AgentRequest{ID: "id1", Endpoint: "endpoint1"},
			AgentRequest{Endpoint: "endpoint2"}, // No ID, should be generated
		}

		results, err := storage.storeAll(StorageTableNameAgentRequests, reqs)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		var parsedResults []*AgentRequest

		for _, result := range results {
			res, err := parseResult[AgentRequest](result)
			assert.NoError(t, err)
			parsedResults = append(parsedResults, res)
		}

		assert.Equal(t, "id1", parsedResults[0].ID)
		assert.Equal(t, "endpoint1", parsedResults[0].Endpoint)
		assert.NotEmpty(t, parsedResults[1].ID)
		assert.Equal(t, "endpoint2", parsedResults[1].Endpoint)
	})

	t.Run("stores multiple items with numeric IDs", func(t *testing.T) {
		type NumericIDStruct struct {
			ID       int    `json:"id,omitempty"`
			Endpoint string `json:"endpoint"`
		}
		reqs := []interface{}{
			NumericIDStruct{ID: 123, Endpoint: "endpoint1"},
			NumericIDStruct{Endpoint: "endpoint2"}, // No ID, should generate numeric
		}

		results, err := storage.storeAll(StorageTableNameAgentRequests, reqs)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		// Check first item with provided numeric ID
		res1, ok := results[0].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, 123, res1["id"])
		assert.Equal(t, "endpoint1", res1["endpoint"].(string))

		// Check second item with generated numeric ID
		res2, ok := results[1].(map[string]interface{})
		assert.True(t, ok)
		id2, ok := res2["id"].(int)
		assert.True(t, ok)
		assert.Greater(t, id2, 0)
		assert.LessOrEqual(t, id2, 1000000)
		assert.Equal(t, "endpoint2", res2["endpoint"].(string))
	})
}

func TestMemoryStorage_GetAll(t *testing.T) {
	storage := NewMemoryStorage()

	// Store some test data
	reqs := []AgentRequest{
		{ID: "id1", Endpoint: "endpoint1"},
		{ID: "id2", Endpoint: "endpoint2"},
		{ID: "id3", Endpoint: "endpoint1"},
	}

	for _, req := range reqs {
		_, err := storage.store(req.TableName(), req)
		assert.NoError(t, err)
	}

	t.Run("retrieves all matching records", func(t *testing.T) {
		results, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{
			"endpoint": "endpoint1",
		})
		assert.NoError(t, err)
		assert.Len(t, results, 2) // Should find two records with endpoint1

		for _, result := range results {
			res, err := parseResult[AgentRequest](result)
			assert.NoError(t, err)
			assert.Equal(t, "endpoint1", res.Endpoint)
		}
	})

	t.Run("returns empty slice for no matches", func(t *testing.T) {
		results, err := storage.getAll(StorageTableNameAgentRequests, map[string]string{
			"endpoint": "non-existent",
		})
		assert.NoError(t, err)
		assert.Empty(t, results)
	})
}
