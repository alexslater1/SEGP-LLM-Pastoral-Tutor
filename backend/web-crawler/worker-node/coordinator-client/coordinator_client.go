package coordinator_client

import (
	"context"
	"encoding/json"
	"time"
)

type Task struct {
	ID        string      `json:"id"`
	CreatedBy string      `json:"created_by"`
	Params    interface{} `json:"params"`
}

func NewTask(id string, createdBy string, params interface{}) *Task {
	return &Task{
		ID:        id,
		CreatedBy: createdBy,
		Params:    params,
	}
}

func (t *Task) toString() (string, error) {
	json, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

type CoordinatorClientTaskTopic string

func (c CoordinatorClientTaskTopic) String() string {
	return string(c)
}

func (c CoordinatorClientTaskTopic) ProcessingTopicString() string {
	return "processing_" + string(c)
}

const (
	CoordinatorClientTaskTopicUrls CoordinatorClientTaskTopic = "urls"
)

var (
	ErrNoTasksToComplete = &CoordinatorClientNoTasksToComplete{}
	ErrNoTasksCompleted  = &CoordinatorClientNoTasksCompleted{}
)

type CoordinatorClient interface {
	CreateTask(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error
	GetTask(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error)
	GetTaskAndSetProcessing(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error)
	SetProcessed(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error
}

type CoordinatorClientNoTasksToComplete struct {
}

func (r *CoordinatorClientNoTasksToComplete) Error() string {
	return "No tasks to complete"
}

type CoordinatorClientNoTasksCompleted struct {
}

func (r *CoordinatorClientNoTasksCompleted) Error() string {
	return "No tasks completed"
}
