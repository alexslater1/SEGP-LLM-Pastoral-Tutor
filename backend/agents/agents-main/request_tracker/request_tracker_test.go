package request_tracker

import (
	"testing"

	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	store := storage.NewMemoryStorage()
	rt := NewRequestTracker(store)

	id, err := rt.NewRequest("chat_completion")
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}
