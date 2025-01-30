package coordinator_client

import (
	"context"
	"encoding/json"
	"time"
)

type Task struct {
	ID        string                 `json:"id"`
	CreatedBy string                 `json:"created_by"`
	Params    map[string]interface{} `json:"params"`
}

type StoredError struct {
	Error   string                     `json:"error"`
	Task    *Task                      `json:"task"`
	Topic   CoordinatorClientTaskTopic `json:"topic"`
	Created time.Time                  `json:"created"`
}

func CastParams[T any](params map[string]interface{}) (*T, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var paramsMap map[string]interface{}
	if err := json.Unmarshal(jsonParams, &paramsMap); err != nil {
		return nil, err
	}

	var paramsT T
	if err := json.Unmarshal(jsonParams, &paramsT); err != nil {
		return nil, err
	}

	return &paramsT, nil
}

func NewTask(id string, createdBy string, params interface{}) (*Task, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	var paramsMap map[string]interface{}
	if err := json.Unmarshal(jsonParams, &paramsMap); err != nil {
		return nil, err
	}

	return &Task{
		ID:        id,
		CreatedBy: createdBy,
		Params:    paramsMap,
	}, nil
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
	CoordinatorClientTaskTopicRag  CoordinatorClientTaskTopic = "rag"
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

	StoreError(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task, err error) error
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
