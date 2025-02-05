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
	task, err := NewTask("test-id", "test-data", nil)
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
	task2, err := NewTask("test-id-2", "test-data-2", nil)
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

func TestMockCoordinatorClient_CreateTasks(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	params1 := map[string]string{
		"url": "https://ethanhosier.com",
	}
	task1, err := NewTask("test-id-1", "test-data-1", params1)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	params2 := map[string]string{
		"url": "https://ethanhosier.com/blog",
	}
	task2, err := NewTask("test-id-2", "test-data-2", params2)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	client.CreateTasks(ctx, CoordinatorClientTaskTopicUrls, []*Task{task1, task2})

	assert.Equal(t, len(client.tasks[CoordinatorClientTaskTopicUrls.String()]), 2)
}

func TestMockCoordinatorClient_NumTasks(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	numTasks, err := client.NumTasks(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of tasks: %v", err)
	}

	assert.Equal(t, numTasks, 0)

	task1, err := NewTask("test-id-1", "test-data-1", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task1)

	numTasks, err = client.NumTasks(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of tasks: %v", err)
	}

	assert.Equal(t, numTasks, 1)
}

func TestMockCoordinatorClient_NumProcessingTasks(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	numProcessingTasks, err := client.NumProcessingTasks(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of processing tasks: %v", err)
	}

	assert.Equal(t, numProcessingTasks, 0)

	task1, err := NewTask("test-id-1", "test-data-1", nil)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if err := client.CreateTask(ctx, CoordinatorClientTaskTopicUrls, task1); err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	numProcessingTasks, err = client.NumProcessingTasks(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of processing tasks: %v", err)
	}

	assert.Equal(t, numProcessingTasks, 0)

	client.GetTaskAndSetProcessing(ctx, 1*time.Second, CoordinatorClientTaskTopicUrls)

	numProcessingTasks, err = client.NumProcessingTasks(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get number of processing tasks: %v", err)
	}

	assert.Equal(t, numProcessingTasks, 1)
}

func TestMockCoordinatorClient_GetErrors(t *testing.T) {
	client := NewMockCoordinatorClient()
	ctx := context.Background()

	if err := client.StoreError(ctx, CoordinatorClientTaskTopicUrls, nil, fmt.Errorf("an error")); err != nil {
		t.Fatalf("Failed to store error: %v", err)
	}

	if err := client.StoreError(ctx, CoordinatorClientTaskTopicUrls, nil, fmt.Errorf("an error 2")); err != nil {
		t.Fatalf("Failed to store error: %v", err)
	}

	errors, err := client.GetErrors(ctx, CoordinatorClientTaskTopicUrls)
	if err != nil {
		t.Fatalf("Failed to get errors: %v", err)
	}

	assert.Equal(t, len(errors), 2)

	assert.Equal(t, errors[0].Error, "an error")
	assert.Equal(t, errors[1].Error, "an error 2")
}
