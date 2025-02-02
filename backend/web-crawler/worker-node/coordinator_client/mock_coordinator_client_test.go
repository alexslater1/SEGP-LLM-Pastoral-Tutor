package coordinator_client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMockCoordinatorClient(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	// Test CreateTask
	task, err := NewTask("test-id", "test-data", "test-type")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test GetTask
	retrievedTask, err := client.GetTask(ctx, 1*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}
	if retrievedTask.ID != task.ID {
		t.Errorf("Expected task ID %s, got %s", task.ID, retrievedTask.ID)
	}

	// Test GetTaskAndSetProcessing
	task2, err := NewTask("test-id-2", "test-data-2", "test-type-2")
	if err != nil {
		t.Fatalf("Failed to create second task: %v", err)
	}

	err = client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task2)
	if err != nil {
		t.Fatalf("Failed to create second task: %v", err)
	}

	processingTask, err := client.GetTaskAndSetProcessing(ctx, 1*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task and set processing: %v", err)
	}
	if processingTask.ID != task2.ID {
		t.Errorf("Expected task ID %s, got %s", task2.ID, processingTask.ID)
	}

	// Test SetProcessed
	err = client.SetProcessed(ctx, CoordinatorClientTaskTopicUrls, processingTask)
	if err != nil {
		t.Fatalf("Failed to set task as processed: %v", err)
	}

	// Test error cases
	_, err = client.GetTask(ctx, 1*time.Second, CoordinatorClientTaskTopicUrls)
	if err != ErrNoTasksToComplete {
		t.Errorf("Expected ErrNoTasksToComplete, got %v", err)
	}

	err = client.SetProcessed(ctx, CoordinatorClientTaskTopicUrls, task)
	if err != ErrNoTasksCompleted {
		t.Errorf("Expected ErrNoTasksCompleted, got %v", err)
	}
}

func TestMockCoordinatorClient_CreateTask(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	type testParams struct {
		Number int    `json:"number"`
		Name   string `json:"name"`
	}

	task, err := NewTask("test-id", "test-data", testParams{Number: 1, Name: "a name"})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	createdTask, err := client.GetTask(ctx, 1*time.Second, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get task: %v", err)
	}

	parsedParams, err := CastParams[testParams](createdTask.Params)
	if err != nil {
		t.Fatalf("Failed to parse params: %v", err)
	}

	assert.Equal(t, parsedParams.Number, 1)
	assert.Equal(t, parsedParams.Name, "a name")

	t.Logf("createdTask: %+v", createdTask)
}

func TestMockCoordinatorClient_StoreError(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	params := map[string]string{
		"url": "https://ethanhosier.com",
	}

	task, err := NewTask("test-id", "test-data", params)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	client.StoreError(ctx, CoordinatorClientTaskTopicUrls, task, fmt.Errorf("an error"))

	assert.Equal(t, len(client.errors), 1)
	t.Logf("errors: %+v", client.errors)
}
