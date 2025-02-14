package email

import (
	"net/mail"
)

const (
	emailFrom          = "Acme <onboarding@resend.dev>"
	PersonalTutorEmail = "ethanjhosier@gmail.com"
)

type EmailClient interface {
	SendEmail(to string, subject string, htmlBody string) (string, error)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
