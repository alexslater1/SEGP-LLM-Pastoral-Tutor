package agent

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/memory"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func TestReActAgent(t *testing.T) {
	var (
		toolHandler = tools.NewGoogleSearchToolHandler(googleSearch.NewRodClient())
		llm         = llm.NewDeepSeekLLM(os.Getenv("DEEPSEEK_API_KEY"))
		memory      = memory.NewReActMemory()
		storage     = storage.NewLocalStorage()
		knowledge   = knowledge.NewLocalKnowledge()
		agent       = NewReActAgent("You are a ReAct agent", toolHandler, llm, memory, storage, knowledge)
	)

	response, reasoning, err := agent.Run("Who is the headmaster of the secondary school which Dillan Scott (imperial college london) attended?")
	if err != nil {
		t.Fatalf("Error running agent: %v", err)
	}

	t.Logf("Response: %s", *response)
	t.Logf("Reasoning: %s", *reasoning)
}
