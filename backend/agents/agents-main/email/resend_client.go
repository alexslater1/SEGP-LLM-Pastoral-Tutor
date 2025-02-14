package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendClient struct {
	client *resend.Client
}

func NewResendClient(apiKey string) *ResendClient {
	return &ResendClient{client: resend.NewClient(apiKey)}
}

func (c *ResendClient) SendEmail(to string, subject string, htmlBody string) (string, error) {
	if !isValidEmail(to) {
		return "", fmt.Errorf("invalid email address: %s", to)
	}

	params := &resend.SendEmailRequest{
		From:    emailFrom,
		To:      []string{to},
		Html:    htmlBody,
		Subject: subject,
	}

	resp, err := c.client.Emails.Send(params)
	return resp.Id, err
}
