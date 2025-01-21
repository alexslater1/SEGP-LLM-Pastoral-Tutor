package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestCompletions(t *testing.T) {
	nextStep, err := handleNextStep("How can i travel to uni")
	if err != nil {
		t.Errorf("Error handling next step: %v", err)
		return
	}

	t.Logf("Next step: %+v", nextStep)
}
