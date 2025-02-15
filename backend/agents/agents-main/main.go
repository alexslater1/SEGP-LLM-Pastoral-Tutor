package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/agent"
	"github.com/segp/agents-main/api"
	"github.com/segp/agents-main/history"
	"github.com/segp/agents-main/storage"
	"github.com/segp/agents-main/utils"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	var (
		store   = storage.NewSupabaseStorage(utils.Required(os.Getenv("SUPABASE_URL"), "SUPABASE_URL"), utils.Required(os.Getenv("SUPABASE_SERVICE_KEY"), "SUPABASE_SERVICE_KEY"))
		agent   = agent.NewDefaultEventStoringLoggingUserQueryAgent()
		history = history.NewAgentEventHistory(store)
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	listenAddr := flag.String("listen", ":"+port, "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr, store, agent, history)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
