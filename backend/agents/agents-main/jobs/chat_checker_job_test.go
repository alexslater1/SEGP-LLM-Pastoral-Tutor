package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Load .env file before running tests
	err := godotenv.Load("../../../.env")
	if err != nil {
		// Don't fail if .env file is not found, just log it
		println("Warning: .env file not found")
	}

	os.Exit(m.Run())
}

func TestGetStaleSessions(t *testing.T) {
	// Create a new MemoryStorage instance
	memStorage := storage.NewMemoryStorage()
	chatCheckerJob := &ChatCheckerJob{
		store:          memStorage,
		staleWindow:    5 * time.Minute,
	}

	now := time.Now()
	req1Time := now.Add(-12 * time.Minute)
	req2Time := now.Add(-10 * time.Minute)
	req3Time := now.Add(-6 * time.Minute)

	// Add mock data to MemoryStorage
	storage.StoreAll(memStorage, storage.RequestSession{
		SessionID: "session1",
		RequestID: "request1",
		CreatedAt: &req1Time,
	}, storage.RequestSession{
		SessionID: "session1",
		RequestID: "request2",
		CreatedAt: &req2Time,
	}, storage.RequestSession{
		SessionID: "session1",
		RequestID: "request3",
		CreatedAt: &req3Time,
	})

	req4Time := now.Add(-50 * time.Minute)
	req5Time := now.Add(-20 * time.Minute)
	req6Time := now.Add(-10 * time.Minute)

	storage.StoreAll(memStorage, storage.RequestSession{
		SessionID: "session2",
		RequestID: "request4",
		CreatedAt: &req4Time,
	}, storage.RequestSession{
		SessionID: "session2",
		RequestID: "request5",
		CreatedAt: &req5Time,
	}, storage.RequestSession{
		SessionID: "session2",
		RequestID: "request6",
		CreatedAt: &req6Time,
	})

	// Store last checked times for sessions
	nowChecked := now.Add(-10 * time.Minute)
	storage.Store(memStorage, storage.LastCheck{
		ID: "session1",
		CheckedAt: &nowChecked,
	})
	storage.Store(memStorage, storage.LastCheck{
		ID: "session2",
		CheckedAt: &nowChecked,
	})

	// Call the method
	staleSessions, err := chatCheckerJob.getStaleSessions()
	assert.NoError(t, err)

	fmt.Printf("stale %+v", staleSessions)

	// Assertions
	assert.Len(t, staleSessions, 1) // Expecting 1 stale session
	assert.Contains(t, staleSessions, "session1") // Check that the stale session ID is correct

	// Test case 1: No sessions
	memStorage2 := storage.NewMemoryStorage()
	chatCheckerJob2 := &ChatCheckerJob{
		store:          memStorage2,
		staleWindow:    5 * time.Minute,
	}
	staleSessions, err = chatCheckerJob2.getStaleSessions()
	assert.NoError(t, err)
	assert.Len(t, staleSessions, 0) // Expecting no stale sessions

	// Test case 2: All sessions are fresh
	reqFreshTime := now.Add(-1 * time.Minute)
	storage.Store(memStorage2, storage.RequestSession{
		SessionID: "session3",
		RequestID: "request7",
		CreatedAt: &reqFreshTime,
	})
	storage.Store(memStorage2, storage.LastCheck{
		ID: "session3",
		CheckedAt: &now, // Last checked time is now
	})
	staleSessions, err = chatCheckerJob2.getStaleSessions()
	assert.NoError(t, err)
	assert.Len(t, staleSessions, 0) // Expecting no stale sessions

	// Test case 3: Multiple stale sessions
	reqStaleTime1 := now.Add(-9 * time.Minute)
	reqStaleTime2 := now.Add(-6 * time.Minute)
	storage.StoreAll(memStorage2, storage.RequestSession{
		SessionID: "session4",
		RequestID: "request8",
		CreatedAt: &reqStaleTime1,
	}, storage.RequestSession{
		SessionID: "session5",
		RequestID: "request9",
		CreatedAt: &reqStaleTime2,
	})
	storage.Store(memStorage2, storage.LastCheck{
		ID: "session4",
		CheckedAt: &nowChecked, // Set last checked time
	})
	storage.Store(memStorage2, storage.LastCheck{
		ID: "session5",
		CheckedAt: &nowChecked, // Set last checked time
	})
	staleSessions, err = chatCheckerJob2.getStaleSessions()
	assert.NoError(t, err)
	assert.Len(t, staleSessions, 2) // Expecting 2 stale sessions
	assert.Contains(t, staleSessions, "session4")
	assert.Contains(t, staleSessions, "session5")

	// Test case 4: Edge case with boundary conditions
	reqBoundaryTime1 := now.Add(-10 * time.Minute)
	reqBoundaryTime2 := now.Add(-5 * time.Minute)
	storage.StoreAll(memStorage2, storage.RequestSession{
		SessionID: "session6",
		RequestID: "request10",
		CreatedAt: &reqBoundaryTime1,
	},
	storage.RequestSession{
		SessionID: "session7",
		RequestID: "request11",
		CreatedAt: &reqBoundaryTime2,
	})
	storage.Store(memStorage2, storage.LastCheck{
		ID: "session6",
		CheckedAt: &nowChecked, // Set last checked time
	})
	storage.Store(memStorage2, storage.LastCheck{
		ID: "session7",
		CheckedAt: &nowChecked, // Set last checked time
	})
	staleSessions, err = chatCheckerJob2.getStaleSessions()
	assert.NoError(t, err)
	// Expect 3 stale sessions, as session6 is exactly on the last window boundary and session7 has just gone stale
	assert.Len(t, staleSessions, 3) 
	assert.Contains(t, staleSessions, "session4")
	assert.Contains(t, staleSessions, "session5")
	assert.Contains(t, staleSessions, "session7")
}

func TestProcessChat(t *testing.T) {
	chatCheckerJob := &ChatCheckerJob{
		llm: 			llm.NewGeminiLLM(context.Background(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY is not set")),
	}

	messages1 := []string{
    `Query: I have three assignments due next week, and I don't know how I'm going to finish them all.
	Response: That sounds overwhelming. Have you spoken to your professors about extensions or support options?`,
    `Query: No, I don't think they'd care. Everyone else seems to be managing just fine, but I'm falling apart.
	Response: It's okay to ask for help! The university has academic support services - would you like me to connect you?`,
    `Query: Maybe... I just feel like such a failure.
	Response: You're not a failure! Many students struggle with workload. Let me provide some resources that might help.`,
	}

	res1, err1 := chatCheckerJob.processChat(messages1)
	if err1 != nil {
		log.Printf("Error processing chat: %v", err1)
	}

	log.Printf("CHAT1 Flags: %+v", res1)

	messages2 := []string{
    `Query: Hi!
	Response: Hello there! How are you?`,
    `Query: Not great to be honest.
	Response: Oh, that's not good. Anything I can do to help?`,
	}

	res2, err2 := chatCheckerJob.processChat(messages2)
	if err2 != nil {
		log.Printf("Error processing chat: %v", err2)
	}

	log.Printf("CHAT2 Flags: %+v", res2)
}

func TestRun(t *testing.T) {
	chatCheckerJob := &ChatCheckerJob{
		store:       storage.NewMemoryStorage(),
		staleWindow: 5 * time.Minute,
		history:     history.NewLocalHistory(),
		llm:        llm.NewMockLLM(),
	}

	now := time.Now()
	req1Time := now.Add(-6 * time.Minute)
	req2Time := now.Add(-8 * time.Minute)

	// Add mock data to MemoryStorage
	storage.StoreAll(chatCheckerJob.store, storage.RequestSession{
		SessionID: "session1",
		RequestID: "request1",
		CreatedAt: &req1Time,
	}, storage.RequestSession{
		SessionID: "session1",
		RequestID: "request2",
		CreatedAt: &req2Time,
	})

	// Store last checked times for sessions
	nowChecked := now.Add(-10 * time.Minute)
	storage.Store(chatCheckerJob.store, storage.LastCheck{
		ID:        "session1",
		CheckedAt: &nowChecked,
	})

	// Mock the message history for the session using LocalHistory
	chatHistory := []string{
		"Query: Hi!",
		"Response: Hello there! How are you?",
		"Query: Not great to be honest.",
		"Response: Oh, that's not good. Anything I can do to help?",
	}
	chatCheckerJob.history.(*history.LocalHistory).AddMessageHistory(chatHistory)

    // Add mock LLM response
	chatCheckerJob.llm.(*llm.MockLLM).NewCallChain().ThenStructured(`[{"Flag": "Mental health concerns", "Reason": "The student stated they were 'Not great to be honest', indicating potential emotional distress. The conversation ended abruptly after the AI offered help, suggesting the student may have disengaged due to these concerns."}]`).Set()

	// Run the job
	chatCheckerJob.Run()

	// Verify the expected output
	// Check if the last check was updated
	lastCheckData, err := storage.Get[storage.LastCheck](chatCheckerJob.store, "session1")
	assert.NoError(t, err)
	assert.NotNil(t, lastCheckData)
	assert.True(t, lastCheckData.CheckedAt.After(nowChecked), "Last check time should be updated")
}
