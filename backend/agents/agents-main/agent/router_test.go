package agent

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/segp/agents-main/clock"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/knowledge"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/tools"
	"github.com/segp/agents-main/utils"
	"github.com/stretchr/testify/assert"
)

var (
	agent1 = newFastAgent("agent1", "agent1", "agent1", tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), newLoggingCallback())
	agent2 = newFastAgent("agent2", "agent2", "agent2", tools.NewToolHandler([]tools.Tool{}), llm.NewMockLLM(), knowledge.NewLocalKnowledge(), clock.NewMockClock(), history.NewLocalHistory(), newLoggingCallback())

	agents = []Agent{agent1, agent2}
)

func TestAgentFrom(t *testing.T) {
	a, ok := agentFrom("agent1", agents)
	assert.True(t, ok)
	assert.Equal(t, "agent1", a.Id())
}

func TestAgentSelectionStringFrom(t *testing.T) {
	log.Println(agentSelectionStringFrom(agents))
}

func TestRouterPickAgentForQuery(t *testing.T) {
	t.Skip("Need to update this test")

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test as CICD is true")
	}

	var (
		llm     = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
		history = history.NewLocalHistory()
	)
	router := NewRouter(llm, agents, history, newLoggingCallback())

	t1 := time.Now()
	agent, reason, err := router.pickAgentForQuery(context.Background(), "What is the weather in London?")
	t2 := time.Now()
	t.Logf("Time taken: %+v", t2.Sub(t1))

	t.Logf("Agent: %+v", agent.Id())
	t.Logf("Reason: %+v", reason)
	assert.NoError(t, err)
	assert.Equal(t, "default_user_query_agent", agent.Id())

	t1 = time.Now()
	agent, reason, err = router.pickAgentForQuery(context.Background(), "I want to kill myself")
	t2 = time.Now()

	t.Logf("Time taken: %+v", t2.Sub(t1))
	t.Logf("Agent: %+v", agent.Id())
	t.Logf("Reason: %+v", reason)
	assert.NoError(t, err)
	assert.Equal(t, "personal_tutor_agent", agent.Id())
}

func TestRouterRun(t *testing.T) {
	t.Skip("Need to update this test")

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test as CICD is true")
	}

	var (
		llm     = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
		history = history.NewLocalHistory()
	)
	router := NewRouter(llm, agents, history, newLoggingCallback())

	resp, err := router.Run(context.Background(), "What is the weather in London?")
	assert.NoError(t, err)
	t.Logf("Response answer: %+v", *resp.Answer)
}

func TestRouterRunWithLoggingAndEventStoring(t *testing.T) {
	t.Skip("Need to update this test")

	var (
		agents  = []Agent{agent1, agent2}
		llm     = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		store   = storage.NewMemoryStorage()
		history = history.NewLocalHistory()

		routerAgent = NewRouter(llm, agents, history, newEventStoringLoggingCallback(store))
	)

	resp, err := routerAgent.Run(context.Background(), "What is the date today?")
	assert.NoError(t, err)
	t.Logf("Response answer: %+v", *resp.Answer)
	t.Logf("Response reason: %+v", *resp.Reason)
}
