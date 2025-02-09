package main

import (
	"flag"
	"log"

	"github.com/ethanhosier/web-crawler-coordinator/api"
	"github.com/joho/godotenv"
)

// docker run -p 8080:8080 -e REDIS_ADDRESS=host.docker.internal:6379 coordinator

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	listenAddr := flag.String("listen", ":80", "HTTP server listen address")
	flag.Parse()

	server := api.NewServer(*listenAddr)
	log.Printf("Starting server on %s", *listenAddr)
	log.Fatal(server.Start())
}
