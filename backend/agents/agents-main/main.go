package main

import (
	"context"
	"flag"
	"time"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api"
	"github.com/segp/agents-main/clock"
	googleSearch "github.com/segp/agents-main/google_search"
	"github.com/segp/agents-main/imperial_apis"
	"github.com/segp/agents-main/knowledge"

	// "github.com/segp/agents-main/email"
	"log"
	"os"

	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/jobs"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
	// "time"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var (
		store   = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE_KEY"))
		history = history.NewAgentEventHistory(store)

		llm = llm.NewGeminiLLM(context.TODO(), os.Getenv("GEMINI_API_KEY"))

		googleSearchClient = googleSearch.NewRodClient()
		searchKnowledge    = knowledge.NewRAGKnowledge(os.Getenv("RAG_BASE_URL"))

		imperialApiHandler = imperial_apis.NewDefaultImperialApiHandler()

		agentProvider = agent.NewAgentProvider(store, llm, clock.NewRealClock(), history, googleSearchClient, searchKnowledge, imperialApiHandler)
		js            = []jobs.Job{
			// jobs.NewChatCheckerJob(store, history, llm, email.NewResendClient(utils.Required(os.Getenv("RESEND_API_KEY"), "RESEND_API_KEY")), 10*time.Second),
			jobs.NewAdjustPromptsJob(store, history, llm, agentProvider, 10*time.Second),
		}
		jobManager = jobs.NewJobManager(js)
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenAddr := flag.String("listen", ":"+port, "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr, store, history, jobManager, llm, agentProvider)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
