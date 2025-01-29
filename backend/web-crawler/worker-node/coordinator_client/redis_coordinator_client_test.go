package coordinator_client

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

	taskParams := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("2b665be2-80b7-40d4-9117-a6e9794afe97", "asdasdasd", taskParams)
	if err != nil {
		t.Fatalf("Failed to    create ta sk:   %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task      : %v", err)
	}
}

func TestRedisCoordinatorClientCreate100Tasks(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	for i := 0; i < 100; i++ {
		params := map[string]string{
			"url": "https://ethanhosier.com",
		}

		task, err := NewTask(uuid.New().String(), "asdasdasd", params)
		if err != nil {
			t.Fatalf("Failed to create  task: %v", err)
		}

		err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
		if err != nil {
			t.Fatalf("Failed to crea te task: %v", err)
		}
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
	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", "test")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD ")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	client.SetProcessed(context.Background(), CoordinatorClientTaskTopicUrls, task)
}

func TestRedisCoordinatorClientCreateGetTask(t *testing.T) {
	if os.Getenv("CICD") == "true" {
		t.Skip("Skipping test in CICD")
	}

	client := NewRedisCoordinatorClient(context.Background(), "localhost:6379", "", 0)

	type TestStruct struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
	}

	params := TestStruct{Number: 1, Name: "a name"}

	task, err := NewTask("37407602-a309-4afd-8b77-efa91d808bf3", "asdasdasd", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.CreateTask(context.Background(), CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task, err = client.GetTask(context.Background(), 5*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	parsedParams, err := CastParams[TestStruct](task.Params)
	if err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}

	assert.Equal(t, parsedParams.Number, params.Number)
	assert.Equal(t, parsedParams.Name, params.Name)
}
