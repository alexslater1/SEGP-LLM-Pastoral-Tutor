package agent

import (
	"context"
	"log"
	"os"
	"strings"

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

	Id() string
	Description() string
}

func handleEvent(callback func(event AgentEvent), event AgentEvent) {
	if callback == nil {
		return
	}

	go callback(event)
}

func NewEventStoringLoggingCallback(store storage.Storage) func(event AgentEvent) {
	return func(event AgentEvent) {
		log.Printf("[Event]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
		_, err := storage.Store(store, storage.NewAgentEvent(event.RequestID, string(event.Type), event.Data))
		if err != nil {
			log.Printf("[EventStoringCallback] Error storing event: %v. Continuing...", err)
		}
	}
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
		NewEventStoringLoggingCallback(supabaseStore),
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
		"You are a personal tutor agent. You are meant to provide support for a student at imperial college london. You are a layer between the students and their personal tutor. Students interact with you via a chatbot. You must also always reply to the user in a way which is supportive. Your response should seem as if it came from a human. If you have an idea of a specific way to steer the user, you should do that. Examples would include offering to draft an email, offering to do more research on a specific topic, offering to do a task for the user, etc. Whatever would be most helpful for what the user has asked.",

		toolHandler,
		geminiLlm,
		knowledge,
		realClock,
		agentEventHistory,
		NewEventStoringLoggingCallback(supabaseStore),
	)
}

// ^ FINISH IGNORING
func newSpecializedAgent(id string, prompt string, description string, apiFuncs ...interface{}) *FastAgent {
	var (
		basePrompt = ""

		geminiLlm      = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
		ragKnowledge   = knowledge.NewRAGKnowledge(utils.Required(os.Getenv("RAG_BASE_URL"), "RAG_BASE_URL"))
		extraKnowledge = knowledge.NewExtraKnowledge(apiFuncs...)
		// conjoinedKnowledge = knowledge.NewConjoinedKnowledge(ragKnowledge, extraKnowledge)
		realClock = clock.NewRealClock()

		googleSearchClient = googleSearch.NewRodClient()
		supabaseStore      = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE"))
		agentEventHistory  = history.NewAgentEventHistory(supabaseStore)

		toolHandler = tools.NewToolHandler([]tools.Tool{
			tools.NewSearchTool(ragKnowledge, tools.NewGoogleSearchFirstResultsPageContentsTool(googleSearchClient, 3)),
		})

		callback = NewEventStoringLoggingCallback(supabaseStore)
	)

	finalPrompt := prompt + basePrompt

	return newFastAgent(id, description, finalPrompt, toolHandler, geminiLlm, extraKnowledge, realClock, agentEventHistory, callback)
}
