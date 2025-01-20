package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/api"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	listenAddr := flag.String("listen", ":8080", "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
