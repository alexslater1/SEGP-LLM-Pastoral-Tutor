package agent

import (
	"context"
	"os"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/memory"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

type Agent interface {
	Run(input string, requestId string) (*string, *string, error)
	Subscribe() <-chan AgentEvent
	Unsubscribe(ch <-chan AgentEvent)
}

func NewDefaultReActAgent() Agent {
	var (
		toolHandler = tools.NewDefaultToolHandler(googleSearch.NewNonHeadlessRodClient(), knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL")))
		llm         = llm.NewOpenAiLLM(os.Getenv("OPENAI_API_KEY"))
		memory      = memory.NewReActMemory()
		history     = history.NewLocalHistory()
		knowledge   = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
		agent       = NewReActAgent("You are a ReAct agent", toolHandler, llm, memory, history, knowledge)
	)

	return agent
}

func NewDefaultFastAgent() Agent {
	var (
		googleSearchClient = googleSearch.NewRodClient()
		toolHandler        = tools.NewGoogleSearchToolHandler(googleSearchClient)
		llm                = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		knowledge          = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
	)

	return NewFastAgent(
		"You are a ReAct agent",
		toolHandler,
		llm,
		knowledge,
	)
}

func NewDefaultLoggingReActAgent() Agent {
	return NewLoggingAgent(NewDefaultReActAgent())
}

func NewDefaultEventStoringReActAgent() Agent {
	return NewEventStoringAgent(NewDefaultReActAgent(), storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY")))
}

func NewDefaultEventStoringLoggingReActAgent() Agent {
	return NewLoggingAgent(NewDefaultEventStoringReActAgent())
}

func NewDefaultLoggingFastAgent() Agent {
	return NewLoggingAgent(NewDefaultFastAgent())
}

func NewDefaultEventStoringFastAgent() Agent {
	return NewEventStoringAgent(NewDefaultFastAgent(), storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY")))
}

func NewDefaultEventStoringLoggingFastAgent() Agent {
	return NewLoggingAgent(NewDefaultEventStoringFastAgent())
}
