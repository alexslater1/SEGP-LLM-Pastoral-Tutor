package email

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/utils"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("Error loading .env file")
	}
	os.Exit(m.Run())
}

func TestResendClient(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CI/CD")
	}

	client := NewResendClient(utils.Required(os.Getenv("RESEND_API_KEY"), "RESEND_API_KEY"))
	client.SendEmail("ethanjhosier@gmail.com", "Test Subject", "Test Body")
}
