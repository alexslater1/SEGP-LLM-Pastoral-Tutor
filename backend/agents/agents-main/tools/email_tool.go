package tools

import (
	"fmt"

	"github.com/segp/agents-main/email"
)

type EmailTool struct {
	emailClient  email.EmailClient
	email        string
	instructions string
	to           string
	subject      string
}

func NewEmailTool(to string, subject string, ec email.EmailClient, email string, instructions string) *EmailTool {
	return &EmailTool{
		emailClient:  ec,
		email:        email,
		instructions: instructions,
		to:           to,
		subject:      subject,
	}
}

func (e *EmailTool) Definition() ToolDefinition {
	return ToolDefinition{
		Name:        "email",
		Description: fmt.Sprintf("Send one email to: %s. Instructions: %s. Note that this tool must only be used at most once per query.", e.email, e.instructions),
		Parameters: []Parameter{
			{Name: "html_body", Description: "The body of the email in html format", Type: ParameterTypeString},
		},
	}
}

func (e *EmailTool) SendEmail(htmlBody string) (string, error) {
	return e.emailClient.SendEmail(e.to, e.subject, htmlBody)
}
