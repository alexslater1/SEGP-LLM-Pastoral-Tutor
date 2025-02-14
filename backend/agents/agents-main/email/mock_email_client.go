package email

import "log"

type SentEmail struct {
	EmailFrom string
	EmailTo   string
	Subject   string
	HtmlBody  string
}

type MockEmailClient struct {
	SentEmails []SentEmail
}

func NewMockEmailClient() *MockEmailClient {
	return &MockEmailClient{SentEmails: []SentEmail{}}
}

func (c *MockEmailClient) SendEmail(to string, subject string, htmlBody string) (string, error) {
	c.SentEmails = append(c.SentEmails, SentEmail{EmailFrom: emailFrom, EmailTo: to, Subject: subject, HtmlBody: htmlBody})

	log.Println("Sent mock email to", to, "with subject", subject, "and body", htmlBody)
	return "123", nil
}

func (c *MockEmailClient) GetSentEmails() []SentEmail {
	return c.SentEmails
}
