package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/segp/agents-main/email"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
)

const (
	flagging_prompt = `You are analyzing chat interactions between university students and an AI chatbot. The following chat you are reviewing recently went inactive, meaning the student has stopped responding. Your task is to assess whether the chat contains potential concerns that require follow-up from the pastoral care team. This should only happen in the case where the student is explicitly showing signs of distress or is in need of support. Examples include:
	-	Showing signs of distress
	-	Showing signs of being overwhelmed
	-	Showing signs of being stressed
	-	Showing signs of being anxious
	-	Showing signs of being depressed

	Here is the chat history:
	%+v

	For the given interaction, analyze the messages and the context of the conversation to decide whether an email should be sent to the pastoral care team. If so, in said email, give them relevant details and some context on the situation. 
	
	Notes:
	-	The email should be sent from the perspective of the chatbot (called Amanda)
	-	Do not include a subject
	- The response must be in html format (so structure it correctly)
	- After signing off, you must include the exact sent from the user and the agent in a structured manner`
)

type ChatCheckerJob struct {
	store       storage.Storage
	history     history.History
	llm         llm.LLM
	emailClient email.EmailClient
	staleWindow time.Duration
}

func NewChatCheckerJob(store storage.Storage, history history.History, llm llm.LLM, emailClient email.EmailClient, staleWindow time.Duration) *ChatCheckerJob {
	return &ChatCheckerJob{
		store,
		history,
		llm,
		emailClient,
		staleWindow,
	}
}

func (c *ChatCheckerJob) Run() error {
	staleSessionIDs, err := c.getStaleSessions()
	if err != nil {
		return fmt.Errorf("error getting stale sessions: %v", err)
	}

	fmt.Printf("Processing session ids: %+v\n", staleSessionIDs)

	staleChats := make(map[string][]string)
	for _, sessionID := range staleSessionIDs {
		chats, err := c.history.GetMessageHistory(sessionID)
		if err != nil {
			return fmt.Errorf("error getting stale chats: %v", err)
		}

		staleChats[sessionID] = chats
	}

	for _, chat := range staleChats {
		emailBody, sendEmail, err := c.processChat(chat)
		if err != nil {
			return fmt.Errorf("error processing chat: %v", err)
		}

		if sendEmail {
			log.Print("Sending email")
			c.emailClient.SendEmail(email.PersonalTutorEmail, "URGENT: Student requires you attention", emailBody)
		}
	}

	_, err = storage.Store(c.store, storage.NewChatCheck())
	if err != nil {
		return fmt.Errorf("error updating last check time: %v", err)
	}

	return nil
}

func (c *ChatCheckerJob) processChat(chat []string) (string, bool, error) {

	type EmailStructuredOutput struct {
		SendEmail bool   `json:"send_email"`
		EmailBody string `json:"email_body"`
	}

	structuredEmailResponse, err := c.llm.StructuredOutputCompletion(context.Background(), fmt.Sprintf(flagging_prompt, chat), EmailStructuredOutput{})
	if err != nil {
		return "", false, err
	}

	var parsedStructuredEmailResponse EmailStructuredOutput
	if err := json.Unmarshal([]byte(*structuredEmailResponse), &parsedStructuredEmailResponse); err != nil {
		return "", false, err
	}

	return parsedStructuredEmailResponse.EmailBody, parsedStructuredEmailResponse.SendEmail, nil
}

func (c *ChatCheckerJob) getStaleSessions() ([]string, error) {
	staleThreshold := time.Now().Add(-c.staleWindow)

	requestSessions, err := storage.GetAll[storage.RequestSession](c.store, nil)
	if err != nil {
		return nil, err
	}

	latestRequestSessions := make(map[string]storage.RequestSession)
	for _, requestSession := range requestSessions {
		if existing, found := latestRequestSessions[requestSession.SessionID]; !found || requestSession.CreatedAt.After(*existing.CreatedAt) {
			latestRequestSessions[requestSession.SessionID] = requestSession
		}
	}

	chatCheckData, err := storage.GetAll[storage.ChatCheck](c.store, nil)
	if err != nil {
		return nil, err
	}
	sort.Slice(chatCheckData, func(i, j int) bool {
		return chatCheckData[i].CreatedAt.After(*chatCheckData[j].CreatedAt)
	})

	if len(chatCheckData) == 0 {
		var sessionIDs []string
		for sessionID := range latestRequestSessions {
			sessionIDs = append(sessionIDs, sessionID)
		}
		return sessionIDs, nil
	}

	lastChecked := chatCheckData[0].CreatedAt

	var staleSessionIDs []string
	for _, requestSession := range latestRequestSessions {
		if requestSession.CreatedAt.Before(staleThreshold) && requestSession.CreatedAt.After(lastChecked.Add(-c.staleWindow)) {
			staleSessionIDs = append(staleSessionIDs, requestSession.SessionID)
		}
	}

	return staleSessionIDs, nil
}

func (c *ChatCheckerJob) Interval() time.Duration {
	return c.staleWindow
}

func (c *ChatCheckerJob) Name() string {
	return "ChatCheckerJob"
}
