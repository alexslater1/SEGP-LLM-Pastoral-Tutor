package crew

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/agent"

	"github.com/segp/agents-main/entity"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	os.Exit(m.Run())
}

func TestCrew(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("skipping test in CI")
	}

	var (
		userQueryAgent     = agent.NewDefaultEventStoringLoggingUserQueryAgent()
		personalTutorAgent = agent.NewDefaultEventStoringLoggingPersonalTutorAgent()
	)

	crew := NewCrew(map[entity.Entity][]entity.Entity{
		userQueryAgent:     {personalTutorAgent, entity.UserEntity},
		personalTutorAgent: {userQueryAgent},
	})

	answer, reason, err := crew.Run(context.Background(), "How many children does the current richest man in the world have? Email the answer to my personal tutor", userQueryAgent.Id())
	if err != nil {
		t.Fatalf("error running crew: %v", err)
	}

	t.Logf("answer: %s", *answer)

	if reason != nil {
		t.Logf("reason: %s", *reason)
	} else {
		t.Log("reason: nil")
	}
}
