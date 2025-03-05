package jobs

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/clock"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/imperial_apis"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
	"github.com/stretchr/testify/assert"
)

func TestAdjustPromptsJob(t *testing.T) {
	if os.Getenv("SUPABASE_URL") == "" {
		t.Skip("SUPABASE_URL is not set")
	}

	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	feedback, err := getFeedback(store, time.Now().Add(-999999*time.Hour))
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

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil, nil)

	feedback, err := getFeedback(store, time.Now().Add(-999999*time.Hour))
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

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil, nil)
	feedback, err := getFeedback(store, time.Now().Add(-99*time.Hour))
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

	job := NewAdjustPromptsJob(store, history.NewAgentEventHistory(store), nil, nil)

	feedback, err := getFeedback(store, time.Now().Add(-999999*time.Hour))
	assert.NoError(t, err)

	enrichedFeedback, err := job.getEnrichedFeedback(feedback)
	assert.NoError(t, err)

	prompt := promptFrom(enrichedFeedback[:5], "CURRENT PROMPT")

	fmt.Println(prompt)
}

func TestUpdatePromptsJob(t *testing.T) {
	var (
		store   = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE_KEY"))
		history = history.NewAgentEventHistory(store)

		llm = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))

		googleSearchClient = googleSearch.NewRodClient()
		searchKnowledge    = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))

		imperialApiHandler = imperial_apis.NewDefaultImperialApiHandler()

		agentProvider = agent.NewAgentProvider(store, llm, clock.NewRealClock(), history, googleSearchClient, searchKnowledge, imperialApiHandler)
	)

	job := NewAdjustPromptsJob(store, history, llm, agentProvider)

	err := job.Run(context.TODO())
	assert.NoError(t, err)

	fmt.Println("Doneeeeeeees")
}
