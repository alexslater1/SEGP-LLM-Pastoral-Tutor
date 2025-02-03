package coordinator_client

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// MockCoordinatorClient implements the coordinator client interface using in-memory storage
type MockCoordinatorClient struct {
	tasks      map[string][]string // topic -> tasks
	processing map[string][]string // topic -> processing tasks
	errors     []string
	mutex      sync.Mutex
}

// NewMockCoordinatorClient creates a new mock coordinator client
func NewMockCoordinatorClient() *MockCoordinatorClient {
	return &MockCoordinatorClient{
		tasks:      make(map[string][]string),
		processing: make(map[string][]string),
		errors:     make([]string, 0),
	}
}

func (m *MockCoordinatorClient) CreateTask(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	taskString, err := task.toString()
	if err != nil {
		return err
	}

	if _, exists := m.tasks[topic.String()]; !exists {
		m.tasks[topic.String()] = make([]string, 0)
	}
	m.tasks[topic.String()] = append(m.tasks[topic.String()], taskString)
	return nil
}

func (m *MockCoordinatorClient) GetTask(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	time.Sleep(timeout)

	topicStr := topic.String()
	if len(m.tasks[topicStr]) == 0 {
		return nil, ErrNoTasksToComplete
	}

	// Get first task
	taskString := m.tasks[topicStr][0]
	// Remove it from the slice
	m.tasks[topicStr] = m.tasks[topicStr][1:]

	var task Task
	if err := json.Unmarshal([]byte(taskString), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (m *MockCoordinatorClient) GetTaskAndSetProcessing(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	time.Sleep(timeout)

	topicStr := topic.String()
	if len(m.tasks[topicStr]) == 0 {
		return nil, ErrNoTasksToComplete
	}

	// Get first task
	taskString := m.tasks[topicStr][0]
	// Remove it from the slice
	m.tasks[topicStr] = m.tasks[topicStr][1:]

	// Add to processing
	processingTopic := topic.ProcessingTopicString()
	if _, exists := m.processing[processingTopic]; !exists {
		m.processing[processingTopic] = make([]string, 0)
	}
	m.processing[processingTopic] = append(m.processing[processingTopic], taskString)

	var task Task
	if err := json.Unmarshal([]byte(taskString), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (m *MockCoordinatorClient) SetProcessed(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	taskString, err := task.toString()
	if err != nil {
		return err
	}

	processingTopic := topic.ProcessingTopicString()
	tasks := m.processing[processingTopic]
	found := false

	// Find and remove the task from processing
	for i, t := range tasks {
		if t == taskString {
			m.processing[processingTopic] = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return ErrNoTasksCompleted
	}

	return nil
}

func (m *MockCoordinatorClient) StoreError(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task, err error) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	storedError := StoredError{
		Error:   err.Error(),
		Task:    task,
		Topic:   topic,
		Created: time.Now(),
	}

	b, err := json.Marshal(storedError)
	if err != nil {
		return err
	}

	m.errors = append(m.errors, string(b))

	return nil
}
