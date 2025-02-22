package agent

import (
	"context"
	"os"

	"github.com/segp/agents-main/clock"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
	"github.com/segp/agents-main/utils"
)

const (
	defaultSubscriberBufferSize = 100
)

type AgentResponse struct {
	Answer *string
	Reason *string
}

type Agent interface {
	Run(ctx context.Context, input string) (*AgentResponse, error)
	Subscribe() <-chan AgentEvent
	Unsubscribe(ch <-chan AgentEvent)

	Id() string
	Description() string
}

// IGNORE THIS
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

		toolHandler = tools.NewGiveAnswerToolHandler([]tools.Tool{
			// tools.NewEmailTool("personal.tutor@imperial.ac.uk", "Personal Tutor", email.NewMockEmailClient(), "To be used to send an email to a personal tutor, in case of a concern."),
		})
	)

	return newFastAgent(
		"personal_tutor_agent",
		"A personal tutor agent, for the student. If there is anything concerning in the message regarding the user, use this agent. It will come up with a response tailored to the user's sitution, relative to imperial college london (which is where the student attends).",
		"You are a personal tutor agent. You are meant to provide support for a student at imperial college london. You are a layer between the students and their personal tutor. Students interact with you via a chatbot. You must also always reply to the user in a way which is supportive. Your response should seem as if it came from a human. If you have an idea of a specific way to steer the user, you should do that. Examples would include offering to draft an email, offering to do more research on a specific topic, offering to do a task for the user, etc. Whatever would be most helpful for what the user has asked.",

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

// ^ FINISH IGNORING

func newSpecializedAgent(id string, prompt string, apiFuncs ...interface{}) *FastAgent {
	var (
		basePrompt = ""

		description        = "TODO"
		geminiLlm          = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
		ragKnowledge       = knowledge.NewRAGKnowledge(utils.Required(os.Getenv("RAG_BASE_URL"), "RAG_BASE_URL"))
		extraKnowledge     = knowledge.NewExtraKnowledge(apiFuncs...)
		conjoinedKnowledge = knowledge.NewConjoinedKnowledge(ragKnowledge, extraKnowledge)
		realClock          = clock.NewRealClock()

		googleSearchClient = googleSearch.NewRodClient()
		supabaseStore      = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE"))
		agentEventHistory  = history.NewAgentEventHistory(supabaseStore)

		toolHandler = tools.NewNoToolGoogleSearchToolHandler(googleSearchClient)
	)

	finalPrompt := basePrompt + prompt

	return newFastAgent(id, description, finalPrompt, toolHandler, geminiLlm, conjoinedKnowledge, realClock, agentEventHistory)
}

// 1: Angelo
func NewAcadmemicSupportStudyAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("academic_support_study_agent", prompt, apiFuncs...))
}

// 2: Angelo
func NewAdminUniServicesAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("admin_uni_services_agent", prompt, apiFuncs...))
}

// 3: Anshul
func NewCareerProfessionalDevelopmentAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("career_professional_development_agent", prompt, apiFuncs...))
}

// 4: Anshul
func NewWellbeingMentalHealthPersonalDevelopmentAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("wellbeing_mental_health_personal_development_agent", prompt, apiFuncs...))
}

// 5: Alex
func NewFinancialAccomodationResourceAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("financial_accomodation_resource_agent", prompt, apiFuncs...))
}

// 6: Alex
func NewCampusLifeSocialAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("campus_life_social_agent", prompt, apiFuncs...))
}

// 7: Teo
func NewAccessibilityDisabilityAgent(prompt string, apiFuncs ...interface{}) Agent {
	return NewLoggingAgent(newSpecializedAgent("accessibility_disability_agent", prompt, apiFuncs...))
}

// 8: Teo
func NewTransitionDiversityMiscAgent(prompt string, apiFuncs ...interface{}) Agent {
	a := newSpecializedAgent("transition_diversity_misc_agent", prompt, apiFuncs...)
	return NewLoggingAgent(a)
}
