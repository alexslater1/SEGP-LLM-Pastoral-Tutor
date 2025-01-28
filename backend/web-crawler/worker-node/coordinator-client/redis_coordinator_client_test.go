package coordinator_client

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	os.Exit(m.Run())
}

func TestRedisCoordinatorClientCreateTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task := NewTask("2b665be2-80b7-40d4-9117-a6e9794afe97", "asdasdasd", "test")
	err := client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task : %v", err)
	}
}

func TestRedisCoordinatorClientGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}
	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientGetTaskAndSetProcessing(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	task, err := client.GetTaskAndSetProcessing(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	t.Logf("Task: %+v", task)
}

func TestRedisCoordinatorClientSetProcessed(t *testing.T) {
	task := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", "test")

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD ")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	client.SetProcessed(context.Background(), CoordinatorClientTaskTopicUrls, task)
}
