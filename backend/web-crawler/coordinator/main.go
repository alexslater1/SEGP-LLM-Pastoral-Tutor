package main

import (
	// "context"
	// "fmt"
	"log"
	"os"
	// "os"

	"github.com/ethanhosier/web-crawler-coordinator/aws"
	"github.com/ethanhosier/web-crawler-coordinator/coordinator_client"
	"github.com/ethanhosier/web-crawler-coordinator/utils"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type ScraperWorkerParams struct {
	URL string `json:"url"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisAddress := utils.Required(os.Getenv("REDIS_ADDRESS"), "REDIS_ADDRESS")
	redisPort := utils.Required(os.Getenv("REDIS_PORT"), "REDIS_PORT")
	redisPassword := utils.Required(os.Getenv("REDIS_PASSWORD"), "REDIS_PASSWORD")
	redisDB := utils.RequiredInt(os.Getenv("REDIS_DB"), "REDIS_DB")

	// coordinatorClient := coordinator_client.NewMockCoordinatorClient(context.Background(), fmt.Sprintf("%s:%s", redisAddress, redisPort), redisPassword, redisDB)

	// urlParams := ScraperWorkerParams{
	// 	URL: "https://www.google.com",
	// }

	// urlTask, err := coordinator_client.NewTask(uuid.New().String(), "scraper", urlParams)
	// if err != nil {
	// 	log.Fatalf("Error creating task: %v", err)
	// }

	// coordinatorClient.CreateTask(context.Background(), coordinator_client.CoordinatorClientTaskTopicUrls, urlTask)

	awsClient := aws.NewAwsClient()
	err := awsClient.LaunchWorkerNode(aws.WorkerTypeScraper, &aws.WorkerNodeParams{
		NumberOfNodes: 1,
		Concurrency:   1,
		RedisAddress:  redisAddress,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
		RedisPort:     redisPort,
	})
	if err != nil {
		log.Fatalf("Error launching worker node: %v", err)
	}
}
