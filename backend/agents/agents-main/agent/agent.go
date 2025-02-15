package agent

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/email"
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
	Entity entity.Entity
	Task   string
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

	// change to set can offload to???
	addCanOffloadToEntity(...entity.Entity)
	canOffloadToEntities() []entity.Entity
	clone() Agent
}

var (
	_ = godotenv.Load("../../../.env")

	googleSearchClient = googleSearch.NewRodClient()
	toolHandler        = tools.NewGoogleSearchToolHandler(googleSearchClient)
	geminiLlm          = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
	// knowledge          = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
	localKnowledge    = knowledge.NewLocalKnowledge()
	realClock         = clock.NewRealClock()
	supabaseStore     = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
	agentEventHistory = history.NewAgentEventHistory(supabaseStore)
)

func NewDefaultUserQueryAgent() Agent {

	return newFastAgent(
		"user_query_agent",
		"An agent that receives the user's query from the frontend. Has a plethora of tools to achieve general tasks.",
		"You are a user query agent. You will be given a real user's query which comes directly from the frontend.",

		toolHandler,
		geminiLlm,
		localKnowledge,
		realClock,
		agentEventHistory,

		entity.UserEntity,
	)
}

func NewPersonalTutorAgent() Agent {
	var (
		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewEmailTool("personal.tutor@imperial.ac.uk", "Personal Tutor", email.NewMockEmailClient(), "personal.tutor@imperial.ac.uk", "To be used to send an email to a personal tutor, in case of a concern."),
		})
	)

	return newFastAgent(
		"personal_tutor_agent",
		"A personal tutor agent. This is an agent which is an expert at handling sensitive topics for the user, such as mental health, or anything where the user need support. If this is applicable, use this agent. It has the ability to email the student's personal tutor too to alert them of any flagged concerns.",
		"You are a personal tutor agent. You are meant to provide support for a student at imperial college london. You are a layer between the students and their personal tutor. Students interact with you via a chatbot. In the case where you have flagged something concerning, you must use the email tool to send an email to the personal tutor, raising this concern and your reasons. You must also always reply to the user in a way which is supportive.",

		toolHandler,
		geminiLlm,
		localKnowledge,
		realClock,
		agentEventHistory,
	)
}

func NewDefaultLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultUserQueryAgent())
}

func NewDefaultEventStoringUserQueryAgent() Agent {
	return NewEventStoringAgent(NewDefaultUserQueryAgent(), supabaseStore)
}

func NewDefaultEventStoringLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultEventStoringUserQueryAgent())
}

func NewDefaultEventStoringLoggingPersonalTutorAgent() Agent {
	return NewLoggingAgent(NewEventStoringAgent(NewPersonalTutorAgent(), supabaseStore))
}
