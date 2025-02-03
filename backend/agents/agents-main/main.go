package main

import (
	"flag"
	"log"
    "os"

	"github.com/joho/godotenv"
	"github.com/segp/agents-main/api"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    listenAddr := flag.String("listen", ":" + port, "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	log.Printf("Starting server on http://localhost%s", *listenAddr)
	log.Fatal(server.Start())
}
