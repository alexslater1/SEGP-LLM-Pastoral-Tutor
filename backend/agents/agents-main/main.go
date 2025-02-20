package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api"
	"github.com/segp/agents-main/crew"
	"github.com/segp/agents-main/email"
	"github.com/segp/agents-main/entity"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/jobs"
	"github.com/segp/agents-main/llm"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var (
		store   = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE_KEY"))
		history = history.NewAgentEventHistory(store)

		userQueryAgent     = agent.NewDefaultEventStoringLoggingUserQueryAgent()
		personalTutorAgent = agent.NewDefaultEventStoringLoggingPersonalTutorAgent()

		crew = crew.NewCrew(map[entity.Entity][]entity.Entity{
			userQueryAgent:     {personalTutorAgent, entity.UserEntity},
			personalTutorAgent: {},
		})

		js = []jobs.Job{
			jobs.NewChatCheckerJob(store, history, llm.NewGeminiLLM(context.Background(), os.Getenv("GEMINI_API_KEY")), email.NewMockEmailClient(), 30*time.Minute),
		}

		jobManager = jobs.NewJobManager(js)
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenAddr := flag.String("listen", ":"+port, "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr, store, userQueryAgent, crew, history, jobManager)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
