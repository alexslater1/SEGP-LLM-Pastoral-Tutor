package main

import (
	"flag"
	"log"

	"github.com/ethanhosier/web-crawler-coordinator/api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	listenAddr := flag.String("listen", ":8080", "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	log.Printf("Starting server on %s", *listenAddr)
	log.Fatal(server.Start())
}
