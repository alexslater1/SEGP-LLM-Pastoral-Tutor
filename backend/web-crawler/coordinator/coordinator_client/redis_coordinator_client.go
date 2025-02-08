package coordinator_client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCoordinatorClient struct {
	redisClient *redis.Client
}

func NewRedisCoordinatorClient(ctx context.Context, address string, password string, db int) *RedisCoordinatorClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return &RedisCoordinatorClient{
		redisClient: redisClient,
	}
}

func (r *RedisCoordinatorClient) CreateTask(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error {
	taskString, err := task.toString()
	if err != nil {
		return err
	}

	return r.redisClient.RPush(ctx, topic.String(), taskString).Err()
}

func (r *RedisCoordinatorClient) CreateTasks(ctx context.Context, topic CoordinatorClientTaskTopic, tasks []*Task) error {
	taskStrings := make([]string, 0, len(tasks))
	for _, task := range tasks {
		taskString, err := task.toString()
		if err != nil {
			return err
		}
		taskStrings = append(taskStrings, taskString)
	}

	taskInterfaces := make([]interface{}, len(taskStrings))
	for i, v := range taskStrings {
		taskInterfaces[i] = v
	}
	return r.redisClient.RPush(ctx, topic.String(), taskInterfaces...).Err()
}

func (r *RedisCoordinatorClient) GetTask(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error) {
	result, err := r.redisClient.BLPop(ctx, timeout, topic.String()).Result()
	if err == redis.Nil {
		return nil, ErrNoTasksToComplete
	}

	if err != nil {
		return nil, err
	}

	var task Task
	err = json.Unmarshal([]byte(result[1]), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *RedisCoordinatorClient) GetTaskAndSetProcessing(ctx context.Context, timeout time.Duration, topic CoordinatorClientTaskTopic) (*Task, error) {
	result, err := r.redisClient.BRPopLPush(ctx, topic.String(), topic.ProcessingTopicString(), timeout).Result()
	if err == redis.Nil {
		return nil, ErrNoTasksToComplete
	}

	if err != nil {
		return nil, err
	}

	var task Task
	err = json.Unmarshal([]byte(result), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *RedisCoordinatorClient) SetProcessed(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task) error {
	taskString, err := task.toString()
	if err != nil {
		return err
	}

	cmd := r.redisClient.LRem(ctx, topic.ProcessingTopicString(), 1, taskString)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if cmd.Val() == 0 {
		return ErrNoTasksCompleted
	}

	return nil
}

func (r *RedisCoordinatorClient) StoreError(ctx context.Context, topic CoordinatorClientTaskTopic, task *Task, err error) error {

	storedError := StoredError{
		Error:   err.Error(),
		Task:    task,
		Topic:   topic,
		Created: time.Now(),
	}

	storedErrorString, err := json.Marshal(storedError)
	if err != nil {
		return err
	}

	return r.redisClient.RPush(ctx, "errors", storedErrorString).Err()
}

func (r *RedisCoordinatorClient) NumTasks(ctx context.Context, topic CoordinatorClientTaskTopic) (int, error) {
	numTasks, err := r.redisClient.LLen(ctx, topic.String()).Result()
	if err != nil {
		return 0, err
	}
	return int(numTasks), nil
}

func (r *RedisCoordinatorClient) NumProcessingTasks(ctx context.Context, topic CoordinatorClientTaskTopic) (int, error) {
	numProcessingTasks, err := r.redisClient.LLen(ctx, topic.ProcessingTopicString()).Result()
	if err != nil {
		return 0, err
	}
	return int(numProcessingTasks), nil
}

func (r *RedisCoordinatorClient) GetErrors(ctx context.Context, topic CoordinatorClientTaskTopic) ([]*StoredError, error) {
	errors, err := r.redisClient.LRange(ctx, "errors", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	storedErrors := make([]*StoredError, 0, len(errors))
	for _, error := range errors {
		var storedError StoredError
		err = json.Unmarshal([]byte(error), &storedError)
		if err != nil {
			return nil, err
		}
		storedErrors = append(storedErrors, &storedError)
	}
	return storedErrors, nil
}
