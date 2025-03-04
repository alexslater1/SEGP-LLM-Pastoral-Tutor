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

		callback = newEventStoringLoggingCallback(supabaseStore)
	)

	finalPrompt := prompt + basePrompt

	return newFastAgent(id, description, finalPrompt, toolHandler, geminiLlm, extraKnowledge, realClock, agentEventHistory, callback)
}

func newEventStoringLoggingCallback(store storage.Storage) func(event AgentEvent) {
	return func(event AgentEvent) {
		log.Printf("[Event]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
		_, err := storage.Store(store, storage.NewAgentEvent(event.RequestID, string(event.Type), event.Data))
		if err != nil {
			log.Printf("[EventStoringCallback] Error storing event: %v. Continuing...", err)
		}
	}
}

func newLoggingCallback() func(event AgentEvent) {
	return func(event AgentEvent) {
		log.Printf("[Event]%s (%s):\n%+v\n\n", strings.ToUpper(string(event.Type)), event.RequestID, event.Data)
	}
}
