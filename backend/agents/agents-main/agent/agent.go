package agent

import (
	"os"

	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/memory"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

type Agent interface {
	Run(input string) (*string, *string, error)
}

func NewDefaultReActAgent() Agent {
	var (
		toolHandler = tools.NewDefaultToolHandler(googleSearch.NewRodClient(), knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL")))
		llm         = llm.NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))
		memory      = memory.NewReActMemory()
		storage     = storage.NewLocalStorage()
		knowledge   = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))
		agent       = NewReActAgent("You are a ReAct agent", toolHandler, llm, memory, storage, knowledge)
	)

	return agent
}
