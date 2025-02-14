package agent

import (
	"context"
	"os"

	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/entity"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

const (
	defaultSubscriberBufferSize = 100
)

type OffloadTask struct {
	EntityID string
	Task     string
}

type AgentResponse struct {
	Answer *string
	Reason *string

	OffloadTask *OffloadTask
}

type Agent interface {
	entity.Entity

	Run(ctx context.Context, input string) (*AgentResponse, error)
	Subscribe() <-chan AgentEvent
	Unsubscribe(ch <-chan AgentEvent)
}

func NewDefaultUserQueryAgent() Agent {
	var (
		googleSearchClient = googleSearch.NewRodClient()
		toolHandler        = tools.NewGoogleSearchToolHandler(googleSearchClient)
		llm                = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		// knowledge          = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
		knowledge = knowledge.NewLocalKnowledge()
		clock     = clock.NewRealClock()
		store     = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
		history   = history.NewAgentEventHistory(store)
	)

	return NewFastAgent(
		"user_query_agent",
		"An agent that receives the user's query from the frontend. Has a plethora of tools to achieve general tasks.",
		"You are a user query agent. You will be given a real user's query which comes directly from the frontend.",

		toolHandler,
		llm,
		knowledge,
		clock,
		history,
	)
}

func NewDefaultLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultUserQueryAgent())
}

func NewDefaultEventStoringUserQueryAgent() Agent {
	return NewEventStoringAgent(NewDefaultUserQueryAgent(), storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY")))
}

func NewDefaultEventStoringLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultEventStoringUserQueryAgent())
}
