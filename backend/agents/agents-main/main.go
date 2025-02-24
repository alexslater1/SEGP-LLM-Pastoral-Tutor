package main

import (
	"context"
	"flag"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api"

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

		// TODO: why does this order matter??? -> does the go routine go out of scope...
		routerAgent = agent.NewLoggingAgent(agent.NewEventStoringAgent(agent.NewRouter(llm, []agent.Agent{
			agent.NewAcadmemicSupportStudyAgent(),
			agent.NewAdminUniServicesAgent(),
			agent.NewCareerProfessionalDevelopmentAgent(),
			agent.NewWellbeingMentalHealthPersonalDevelopmentAgent(),
			agent.NewFinancialAccomodationResourceAgent(),
			agent.NewCampusLifeSocialAgent(),
			agent.NewAccessibilityDisabilityAgent(),
			agent.NewTransitionDiversityMiscAgent(),
			agent.NewGeneralPurposeAgent(),
		}, history), store))

		js = []jobs.Job{
			// jobs.NewChatCheckerJob(store, history, llm, email.NewResendClient(utils.Required(os.Getenv("RESEND_API_KEY"), "RESEND_API_KEY")), 10*time.Second),
		}

		jobManager = jobs.NewJobManager(js)
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenAddr := flag.String("listen", ":"+port, "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr, store, routerAgent, history, jobManager, llm)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
