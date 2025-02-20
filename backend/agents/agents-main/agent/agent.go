package agent

import (
	"context"
	"os"

	"github.com/segp/agents-main/clock"
	// "github.com/segp/agents-main/email"
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

func NewDefaultUserQueryAgent() Agent {

	var (
		googleSearchClient = googleSearch.NewRodClient()
		toolHandler        = tools.NewGoogleSearchToolHandler(googleSearchClient)
		geminiLlm          = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		knowledge          = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
		// localKnowledge    = knowledge.NewLocalKnowledge()
		realClock         = clock.NewRealClock()
		supabaseStore     = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
		agentEventHistory = history.NewAgentEventHistory(supabaseStore)
	)

	return newFastAgent(
		"user_query_agent",
		"An agent that receives the user's query from the frontend. Has a plethora of tools to achieve general tasks.",
		"You are a user query agent. You will be given a real user's query which comes directly from the frontend. You are the only agent who is able to actually communicate with the end user, so remember to recall any information given to you by other agents, and use this in your answer. Answer to the user in a way which is condisderate and helpful. If there is anything concerning the user's wellbeing, you must offload this task to the personal tutor agent.",

		toolHandler,
		geminiLlm,
		knowledge,
		realClock,
		agentEventHistory,

		entity.UserEntity,
	)
}

func NewPersonalTutorAgent() Agent {
	var (
		geminiLlm = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		knowledge = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
		// localKnowledge    = knowledge.NewLocalKnowledge()
		realClock         = clock.NewRealClock()
		supabaseStore     = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
		agentEventHistory = history.NewAgentEventHistory(supabaseStore)

		toolHandler = tools.NewToolHandler([]tools.Tool{
			// tools.NewEmailTool("personal.tutor@imperial.ac.uk", "Personal Tutor", email.NewMockEmailClient(), "To be used to send an email to a personal tutor, in case of a concern."),
		})
	)

	return newFastAgent(
		"personal_tutor_agent",
		"A personal tutor agent, for the student. If there is anything concerning in the message regarding the user, use this agent. It will come up with a response tailored to the user's sitution, relative to imperial college london (which is where the student attends).",
		"You are a personal tutor agent. You are meant to provide support for a student at imperial college london. You are a layer between the students and their personal tutor. Students interact with you via a chatbot. In the case where you have flagged something concerning, you must use the email tool to send an email to the personal tutor, raising this concern and your reasons. You must also always reply to the user in a way which is supportive. Do not tell the user if you have sent an email to the personal tutor.",

		toolHandler,
		geminiLlm,
		knowledge,
		realClock,
		agentEventHistory,
	)
}

func NewDefaultLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultUserQueryAgent())
}

func NewDefaultEventStoringUserQueryAgent() Agent {
	var supabaseStore = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
	return NewEventStoringAgent(NewDefaultUserQueryAgent(), supabaseStore)
}

func NewDefaultEventStoringLoggingUserQueryAgent() Agent {
	return NewLoggingAgent(NewDefaultEventStoringUserQueryAgent())
}

func NewDefaultEventStoringLoggingPersonalTutorAgent() Agent {
	var supabaseStore = storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))
	return NewLoggingAgent(NewEventStoringAgent(NewPersonalTutorAgent(), supabaseStore))
}
