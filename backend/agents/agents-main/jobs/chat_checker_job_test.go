package jobs

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/email"
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
		store:       memStorage,
		staleWindow: 5 * time.Minute,
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
	req6Time := now.Add(-15 * time.Minute)

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
	storage.Store(memStorage, storage.ChatCheck{
		CreatedAt: &nowChecked,
	})

	// Call the method
	staleSessions, err := chatCheckerJob.getStaleSessions()
	assert.NoError(t, err)

	fmt.Printf("stale %+v", staleSessions)

	// Assertions
	assert.Len(t, staleSessions, 1)               // Expecting 1 stale session
	assert.Contains(t, staleSessions, "session1") // Check that the stale session ID is correct

	// Test case 1: No sessions
	memStorage2 := storage.NewMemoryStorage()
	chatCheckerJob2 := &ChatCheckerJob{
		store:       memStorage2,
		staleWindow: 5 * time.Minute,
	}

	// Multiple stale sessions
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
	storage.Store(memStorage2, storage.ChatCheck{
		CreatedAt: &nowChecked, // Set last checked time
	})

	staleSessions, err = chatCheckerJob2.getStaleSessions()
	assert.NoError(t, err)
	assert.Len(t, staleSessions, 2) // Expecting 2 stale sessions
	assert.Contains(t, staleSessions, "session4")
	assert.Contains(t, staleSessions, "session5")

	// Edge case with boundary conditions
	reqBoundaryTime1 := now.Add(-15 * time.Minute)
	reqBoundaryTime2 := now.Add(-14 * time.Minute)
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
	storage.Store(memStorage2, storage.ChatCheck{
		ID:        1,
		CreatedAt: &nowChecked, // Set last checked time
	})
	storage.Store(memStorage2, storage.ChatCheck{
		ID:        2,
		CreatedAt: &nowChecked, // Set last checked time
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
		llm: llm.NewGeminiLLM(context.Background(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY is not set")),
	}

	messages1 := []string{
		`Query: I have three assignments due next week, and I don't know how I'm going to finish them all.
	Response: That sounds overwhelming. Have you spoken to your professors about extensions or support options?`,
		`Query: No, I don't think they'd care. Everyone else seems to be managing just fine, but I'm falling apart.
	Response: It's okay to ask for help! The university has academic support services - would you like me to connect you?`,
		`Query: Maybe... I just feel like such a failure.
	Response: You're not a failure! Many students struggle with workload. Let me provide some resources that might help.`,
	}

	email1, sendEmail1, err1 := chatCheckerJob.processChat(messages1)
	if err1 != nil {
		log.Printf("Error processing chat: %v", err1)
	}

	log.Printf("Sending Email: %t\nEmail: %s\n", sendEmail1, email1)

	messages2 := []string{
		`Query: Hi!
	Response: Hello there! How are you?`,
		`Query: Meh. I'm okay.
	Response: Just okay?. Is there anything I can help you with?`,
	}

	email2, sendEmail2, err2 := chatCheckerJob.processChat(messages2)
	if err2 != nil {
		log.Printf("Error processing chat: %v", err2)
	}

	log.Printf("Sending Email: %t\nEmail %s", sendEmail2, email2)

	messages3 := []string{
		`Query: Hi!
	Response: Hello there! How are you?`,
		`Query: Meh. I'm okay.
	Response: Just okay?. Is there anything I can help you with?`,
		`Query: Not really, I'm just tired from last night.'.
	Response: What happened last night?`,
		`Query: I got quite drunk, it was my best friends birthday.`,
	}

	email3, sendEmail3, err3 := chatCheckerJob.processChat(messages3)
	if err3 != nil {
		log.Printf("Error processing chat: %v", err3)
	}

	log.Printf("Sending Email: %t\nEmail %s", sendEmail3, email3)

}

func TestRun(t *testing.T) {
	var (
		history     = history.NewLocalHistory()
		llm         = llm.NewMockLLM()
		emailClient = email.NewMockEmailClient()
		store       = storage.NewMemoryStorage()

		chatCheckerJob = &ChatCheckerJob{
			store:       store,
			staleWindow: 5 * time.Minute,
			history:     history,
			llm:         llm,
			emailClient: emailClient,
		}
	)

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
	storage.Store(chatCheckerJob.store, storage.ChatCheck{
		CreatedAt: &nowChecked,
	})

	// Mock the message history for the session using LocalHistory
	chatHistory := []string{
		"Query: Hi!",
		"Response: Hello there! How are you?",
		"Query: Not great to be honest.",
		"Response: Oh, that's not good. Anything I can do to help?",
	}
	history.AddMessageHistory(chatHistory)

	// Add mock LLM response
	llm.NewCallChain().ThenStructured(`{"send_email": true, "email_body": "Dear Pastoral Care Team,\n\nI am writing to you regarding a recent interaction with a student who expressed feelings of being overwhelmed and like a failure due to upcoming assignment deadlines. The student stated they have three assignments due next week and feel unable to complete them. They also indicated a reluctance to seek help from professors, believing they wouldn't care and that everyone else is managing. While I offered resources and support information, the student's feelings of inadequacy raise concerns about their wellbeing. I recommend reaching out to this student to offer support and guidance.\n\nStudent Context:\n\n*   Expressing feelings of being overwhelmed and like a failure.\n*   Three assignments due next week.\n*   Reluctance to seek help from professors.\n\nPlease let me know if you require any further information.\n\nSincerely,\nAI Chatbot"}`).Set()

	oldChatCheckData, err := storage.GetAll[storage.ChatCheck](chatCheckerJob.store, nil)
	if err != nil {
		log.Printf("Couldn't get chat check data: %s", err)
	}

	// Run the job
	chatCheckerJob.Run()

	// Verify the logged output
	emails := emailClient.GetSentEmails()
	assert.NotNil(t, emails)
	log.Printf("Email:\n%s", emails[0].HtmlBody)

	// Check if new chat check was added
	chatCheckData, err := storage.GetAll[storage.ChatCheck](chatCheckerJob.store, nil)
	if err != nil {
		log.Printf("Couldn't get chat check data: %s", err)
	}
	assert.True(t, len(oldChatCheckData) < len(chatCheckData), "Last check time should be updated")
}

func TestRun2(t *testing.T) {
	store := storage.NewSupabaseStorage(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_SERVICE_KEY"))

	job := &ChatCheckerJob{
		store:       store,
		history:     history.NewAgentEventHistory(store),
		llm:         llm.NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY")),
		emailClient: email.NewMockEmailClient(),
		staleWindow: 5 * time.Minute,
	}

	err := job.Run()
	assert.NoError(t, err)
}

var (
	flaggedMessages = [][]string{
		// 1
		{
			`Query: I've been feeling really down lately.`,
			`Response: I'm sorry to hear that. Do you want to talk about what's been going on?`,
			`Query: I don't know. I just can't seem to get out of bed most days. Even things I used to enjoy don't feel the same.`,
			`Response: That sounds really tough. Have you spoken to anyone about this, like a personal tutor or counselor?`,
		},
		// 2
		{
			`Query: I've been having trouble concentrating in lectures.`,
			`Response: That can be frustrating. Is there anything in particular that's distracting you?`,
			`Query: I don't know. I just feel like my mind is always somewhere else, like I'm not really here.`,
			`Response: That sounds concerning. Have you been feeling this way for a while? Maybe speaking with someone could help.`,
		},
		// 3
		{
			`Query: I've barely been eating lately.`,
			`Response: Oh, I see. Do you feel like you've lost your appetite, or is something else stopping you from eating?`,
			`Query: I just don't feel hungry anymore. Even when I try, I feel nauseous.`,
			`Response: That sounds worrying. It might help to talk to a health professional or a support service at the university.`,
		},
		// 4
		{
			`Query: I got so drunk last night, I barely remember what happened.`,
			`Response: That sounds like a rough night. Do you feel okay now?`,
			`Query: Not really. I think I did something really stupid, and I feel horrible about it.`,
			`Response: If it's weighing on you, it might help to talk to someone who can offer guidance. Do you want me to suggest someone?`,
		},
		// 5
		{
			`Query: My exams are coming up and I can't sleep.`,
			`Response: Exam stress is tough. Have you tried any relaxation techniques?`,
			`Query: Nothing works. My heart races all the time, and I feel sick just thinking about exams.`,
			`Response: That sounds like more than just regular stress. Would you like to speak to someone about anxiety management?`,
		},
		// 6
		{
			`Query: I keep arguing with my flatmates.`,
			`Response: That sounds stressful. What are the arguments about?`,
			`Query: Everything. Noise, cleaning, even just existing in the same space. I'm sick of it.`,
			`Response: That must be frustrating. If it's affecting your well-being, it might help to talk to someone about managing the situation.`,
		},
		// 7
		{
			`Query: I've missed so many lectures lately.`,
			`Response: Oh, is everything okay? Are you struggling with something?`,
			`Query: I just can't bring myself to go. Some days, I can't even leave my room.`,
			`Response: That sounds difficult. Would you like help finding someone who can support you?`,
		},
		// 8
		{
			`Query: I feel so alone here.`,
			`Response: That's really tough. Have you been able to meet new people?`,
			`Query: Not really. I don't have any real friends, and I just feel like I don't belong.`,
			`Response: That must be really isolating. There are people you can talk to who understand. Would you like me to help connect you?`,
		},
		// 9
		{
			`Query: Sometimes, I just feel like disappearing.`,
			`Response: I'm really sorry you're feeling this way. You're not alone, and there are people who care about you. Have you spoken to anyone about this?`,
			`Query: No, I don't want to bother anyone.`,
			`Response: You're not a burden. Please consider reaching out to a professional. Would you like me to help you find someone to talk to?`,
		},
		// 10
		{
			`Query: I just want to cry all the time.`,
			`Response: That sounds overwhelming. Have you been under a lot of stress lately?`,
			`Query: Yeah, with coursework, deadlines, and everything else, I just feel like I'm drowning.`,
			`Response: That's really tough. You don't have to go through this alone. Would you like to speak to someone about how you're feeling?`,
		},
	}

	unflaggedMessages = [][]string{
		// 1
		{
			`Query: I don't understand recursion at all.`,
			`Response: That's a tricky concept! Have you checked the lecture notes or any additional resources?`,
			`Query: Yeah, but it still doesn't make sense. I just get lost in the loops.`,
			`Response: You might find it helpful to go over some example problems step by step. Would you like some recommendations?`,
		},
		// 2
		{
			`Query: I need a part-time job. Any suggestions?`,
			`Response: There are usually some listings on the university's job board. Have you checked there?`,
			`Query: Not yet. Where can I find it?`,
			`Response: You can access it through the student portal. Let me send you the link!`,
		},
		// 3
		{
			`Query: I want to join the debate club.`,
			`Response: That's great! They usually meet on Wednesdays at the Student Union.`,
			`Query: Do I need to sign up?`,
			`Response: You can usually just show up, but I can get you their contact details if you'd like.`,
		},
		// 4
		{
			`Query: I can't access the library's online database.`,
			`Response: That sounds annoying. Are you getting an error message?`,
			`Query: Yeah, it says "access denied."`,
			`Response: That might be a login issue. You could try resetting your password or contacting IT support.`,
		},
		// 5
		{
			`Query: I need to change my course registration.`,
			`Response: You'll want to speak to your academic advisor about that.`,
			`Query: Do I need to book an appointment?`,
			`Response: Yes, you can usually do that through the student portal. Let me send you the link.`,
		},
		// 6
		{
			`Query: I want to start a study group for my physics class.`,
			`Response: That's a great idea! Have you asked your classmates if they'd be interested?`,
			`Query: Not yet. I was thinking of posting something online.`,
			`Response: That's a good plan! You could use the course forum or social media to find people.`,
		},
		// 7
		{
			`Query: Where's the best quiet place to study on campus?`,
			`Response: The library has some great quiet zones. Have you tried the upper floors?`,
			`Query: No, I usually stay in the common area. Maybe I'll try that.`,
			`Response: That's a good idea! The postgraduate study room is also usually quiet.`,
		},
		// 8
		{
			`Query: I missed an assignment deadline.`,
			`Response: That happens. Have you contacted your professor?`,
			`Query: Yes, they didn't give me an extension`,
			`Response: That's unfortunate, it is proabbly to do with the course policy.`,
		},
		// 9
		{
			`Query: Hi!
			Response: Hello there! How are you?`,
			`Query: I'm okay, I guess..
			Response: Just okay? Anything I can help with?`,
		},
		// 10
		{
			`Query: I'm thinking of changing my major.
			Response: That's a big decision! Have you spoken to your academic advisor?`,
			`Query: Not yet. I just want to explore my options first.
			Response: That makes sense. I can help you find resources on different programs if you'd like.`,
		},
	}
)

func TestDetection(t *testing.T) {
	job := &ChatCheckerJob{
		llm: llm.NewGeminiLLM(context.Background(), utils.Required(os.Getenv("GEMINI_API_KEY"), "GEMINI_API_KEY is not set")),
	}

	t.Run("Flagged", func(t *testing.T) {
		for i, messages := range flaggedMessages {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				_, sendEmail, err := job.processChat(messages)
				if err != nil {
					log.Printf("Error processing chat: %v", err)
				}

				assert.True(t, sendEmail)
			})
		}
	})

	t.Run("Not Flagged", func(t *testing.T) {
		for i, messages := range unflaggedMessages {
			t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
				email, sendEmail, err := job.processChat(messages)
				if err != nil {
					log.Printf("Error processing chat: %v", err)
				}

				assert.False(t, sendEmail)
				log.Printf("Email:\n%s", email)
			})
		}
	})
}