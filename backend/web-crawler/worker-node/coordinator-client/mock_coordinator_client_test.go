package coordinator_client

import (
	"context"
	"testing"
	"time"
)

func TestMockCoordinatorClient(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	// Test CreateTask
	task := NewTask("test-id", "test-data", "test-type")
	err := client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task)
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
	task2 := NewTask("test-id-2", "test-data-2", "test-type-2")
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
