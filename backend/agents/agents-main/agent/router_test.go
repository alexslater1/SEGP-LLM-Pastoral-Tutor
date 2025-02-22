package agent

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
	"github.com/stretchr/testify/assert"
)

func TestAgentFrom(t *testing.T) {
	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
	)
	a, ok := agentFrom("personal_tutor_agent", agents)
	assert.True(t, ok)
	assert.Equal(t, "personal_tutor_agent", a.Id())
}

func TestAgentSelectionStringFrom(t *testing.T) {
	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
	)
	log.Println(agentSelectionStringFrom(agents))
}

func TestRouterPickAgentForQuery(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test as CICD is true")
	}

	var (
		llm    = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
	)
	router := NewRouter(llm, agents)

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
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test as CICD is true")
	}

	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
		llm    = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
	)
	router := NewRouter(llm, agents)

	resp, err := router.Run(context.Background(), "What is the weather in London?")
	assert.NoError(t, err)
	t.Logf("Response answer: %+v", *resp.Answer)
}

func TestRouterSubscribe(t *testing.T) {
	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
		llm    = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
	)
	router := NewRouter(llm, agents)

	router.Subscribe()
	router.Unsubscribe(router.Subscribe())
}

func TestRouterEvents(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test as CICD is true")
	}

	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
		llm    = llm.NewGeminiLLM(context.TODO(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY"))
	)

	logginRouter := NewLoggingAgent(NewRouter(llm, agents))

	resp, err := logginRouter.Run(context.Background(), "What is the date today?")
	assert.NoError(t, err)
	t.Logf("Response answer: %+v", *resp.Answer)
	t.Logf("Response reason: %+v", *resp.Reason)
}

func TestRouterRunWithLoggingAndEventStoring(t *testing.T) {
	var (
		agents = []Agent{NewDefaultUserQueryAgent(), NewPersonalTutorAgent()}
		llm    = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))
		store  = storage.NewMemoryStorage()

		routerAgent = NewEventStoringAgent(NewLoggingAgent(NewRouter(llm, agents)), store)
	)

	resp, err := routerAgent.Run(context.Background(), "What is the date today?")
	assert.NoError(t, err)
	t.Logf("Response answer: %+v", *resp.Answer)
	t.Logf("Response reason: %+v", *resp.Reason)
}
