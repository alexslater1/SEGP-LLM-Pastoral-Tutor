package jobs

import (
	"fmt"
	"os"
	"testing"

	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
	"github.com/stretchr/testify/assert"
)

func TestAdjustPromptsJob(t *testing.T) {
	if os.Getenv("SUPABASE_URL") == "" {
		t.Skip("SUPABASE_URL is not set")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	feedback, err := getFeedback(store)
	assert.NoError(t, err)

	fmt.Println(feedback)
}

func TestIndexOfMatchingRequestID(t *testing.T) {
	agentEvents := []history.MessagesAndActions{
		{RequestID: "1"},
		{RequestID: "2"},
		{RequestID: "3"},
	}
	index, err := indexOfMatchingRequestID(agentEvents, "2")
	assert.NoError(t, err)
	assert.Equal(t, 1, index)
}

func TestGetSingleEnrichedFeedbackFrom(t *testing.T) {
	if os.Getenv("SUPABASE_URL") == "" {
		t.Skip("SUPABASE_URL is not set")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil)

	feedback, err := getFeedback(store)
	assert.NoError(t, err)

	enrichedFeedback, err := job.getSingleEnrichedFeedbackFrom(feedback[3])
	assert.NoError(t, err)

	fmt.Printf("%+v\n", enrichedFeedback)
}

func TestGetEnrichedFeedbackFrom(t *testing.T) {
	if os.Getenv("SUPABASE_URL") == "" {
		t.Skip("SUPABASE_URL is not set")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil)
	feedback, err := getFeedback(store)
	assert.NoError(t, err)

	enrichedFeedback, err := job.getEnrichedFeedback(feedback)
	assert.NoError(t, err)

	fmt.Printf("%+v\n", enrichedFeedback)
}

func TestPromptFrom(t *testing.T) {
	if os.Getenv("SUPABASE_URL") == "" {
		t.Skip("SUPABASE_URL is not set")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil)

	feedback, err := getFeedback(store)
	assert.NoError(t, err)

	enrichedFeedback, err := job.getEnrichedFeedback(feedback)
	assert.NoError(t, err)

	prompt := promptFrom(enrichedFeedback[:5], "CURRENT PROMPT")

	fmt.Println(prompt)
}
