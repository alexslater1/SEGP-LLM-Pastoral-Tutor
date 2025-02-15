package email

import (
	"fmt"
	"log/slog"

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

	slog.Info("Sending email", "to", to, "with subject", subject, "and body", htmlBody)

	params := &resend.SendEmailRequest{
		From:    emailFrom,
		To:      []string{to},
		Html:    htmlBody,
		Subject: subject,
	}

	resp, err := c.client.Emails.Send(params)
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
