package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	assert.True(t, isValidEmail("test@test.com"))
	assert.False(t, isValidEmail("test"))
	assert.False(t, isValidEmail("test@"))
	assert.False(t, isValidEmail("@test.com"))
}

func TestMockEmailClient(t *testing.T) {
	client := NewMockEmailClient()
	client.SendEmail("test@test.com", "Test Subject", "Test Body")

	sentEmails := client.GetSentEmails()
	assert.Equal(t, len(sentEmails), 1)
	assert.Equal(t, sentEmails[0].EmailFrom, emailFrom)
	assert.Equal(t, sentEmails[0].EmailTo, "test@test.com")
	assert.Equal(t, sentEmails[0].Subject, "Test Subject")
	assert.Equal(t, sentEmails[0].HtmlBody, "Test Body")
}
