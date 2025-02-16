package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
)

const (
	flagging_prompt = `You are analyzing chat interactions between university students and an AI chatbot. The following chat you are reviewing recently went inactive, meaning the student has stopped responding. Your task is to assess whether the chat contains potential concerns that require follow-up from the pastoral care team.

	Here is the chat history:
	%+v

	For the given interaction, analyze the messages and the context of the conversation. Flag cases where the student may have disengaged during a sensitive or concerning discussion.

	Possible flags are:
	- Mental health concerns
	- Academic difficulties
	- Personal crises
	- Urgent matters
	- None (If the chat appears to have ended for a neutral or trivial reason)`
)

type flag struct {
		Flag string `json:"flag"`
		Reason string `json:"reason"`
}


type ChatCheckerJob struct {
	store              	storage.Storage
	history				history.History
	llm					llm.LLM
	staleWindow		   	time.Duration
}


func NewChatCheckerJob(store storage.Storage, history history.History, llm llm.LLM, staleWindow time.Duration) *ChatCheckerJob {
	return &ChatCheckerJob{
		store,
		history,
		llm,
		staleWindow,
	}
}


func (c *ChatCheckerJob) Run() {
	staleSessionIDs, err := c.getStaleSessions()
	if err != nil {
		log.Printf("Error getting stale sessions: %v", err)
		return
	}

	staleChats := make(map[string][]string)
	for _, sessionID := range staleSessionIDs {
		chats, err := c.history.GetMessageHistory(sessionID)
		if err != nil {
			log.Printf("Error getting stale chats: %v", err)
		}

		staleChats[sessionID] = chats
	}

	for id, chat := range staleChats {
		res, err := c.processChat(chat)
		if err != nil {
			log.Printf("Error processing chat: %v", err)
		}

		now := time.Now()
		_, err = storage.Store(c.store, storage.NewLastCheck(id, &now))
		if err != nil {
			log.Printf("Error updating last check time for session %s: %v", id, err)
		}

		log.Printf("ID: %s\nFlags: %+v", id, res)
	}
}

func (c *ChatCheckerJob) processChat(chat []string) ([]flag, error) {
	flags, err := c.llm.StructuredOutputCompletion(context.Background(), fmt.Sprintf(flagging_prompt, chat), []flag{})
	if err != nil {
		return nil, err
	}

	var parsedFlags []flag
	if err := json.Unmarshal([]byte(*flags), &parsedFlags); err != nil {
		return nil, err
	}

	return parsedFlags, nil
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
	// log.Printf("latest: %+v", latestRequestSessions)

	var staleSessionIDs []string
	for _, requestSession := range latestRequestSessions {
		var lastChecked time.Time

		lastCheckData, err := storage.Get[storage.LastCheck](c.store, requestSession.SessionID)
		if err == nil && lastCheckData != nil {
			lastChecked = *lastCheckData.CheckedAt
		} else {
			lastChecked = time.Now().Add(-48 * time.Hour)
		}
		// log.Printf("ID: %s lastChecked: %s", requestSession.SessionID, lastChecked)

		if requestSession.CreatedAt.Before(staleThreshold) && requestSession.CreatedAt.After(lastChecked) {
			staleSessionIDs = append(staleSessionIDs, requestSession.SessionID)
		}
	}

	return staleSessionIDs, nil
}