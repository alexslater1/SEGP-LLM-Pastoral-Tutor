package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

//    redisAddress := "redis.worker-node-cluster.local"  // ECS service discovery name

type WorkerType string

const (
	WorkerTypeScraper WorkerType = "scraper"
	WorkerTypeRag     WorkerType = "rag"
)

const (
	defaultRegion = "eu-west-2"
)

type AwsClient struct {
	cfg       aws.Config
	ecsClient *ecs.Client
}

type WorkerNodeParams struct {
	NumberOfNodes int
	Concurrency   int
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RedisPort     string
}

func NewAwsClient() *AwsClient {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(defaultRegion),
	)
	if err != nil {
		panic(err)
	}

	return &AwsClient{
		cfg:       cfg,
		ecsClient: ecs.NewFromConfig(cfg),
	}
}

func (a *AwsClient) LaunchWorkerNode(workerType WorkerType, params *WorkerNodeParams) error {
	// Prepare the network configuration
	networkConfig := types.NetworkConfiguration{
		AwsvpcConfiguration: &types.AwsVpcConfiguration{
			Subnets: []string{
				"subnet-05e199e9936541db1", // public subnet 1 from terraform state
				"subnet-08f965812c1d5ae98", // public subnet 2 from terraform state
			},
			SecurityGroups: []string{
				"sg-0c04d2d1be1d7183d", // worker-node security group from terraform state
			},
			AssignPublicIp: types.AssignPublicIpEnabled,
		},
	}

	taskDefinition := taskDefinitionFromWorkerType(workerType)
	containerOverrides, err := containerOverridesFrom(workerType, params)
	if err != nil {
		return err
	}

	// Run the task
	input := &ecs.RunTaskInput{
		Cluster:              aws.String("worker-node-cluster"),
		TaskDefinition:       aws.String(taskDefinition),
		NetworkConfiguration: &networkConfig,
		LaunchType:           types.LaunchTypeFargate,
		Count:                aws.Int32(int32(params.NumberOfNodes)),
		Overrides: &types.TaskOverride{
			ContainerOverrides: containerOverrides,
		},
	}

	output, err := a.ecsClient.RunTask(context.TODO(), input)
	if err != nil {
		return err
	}

	logInfo(output, workerType, params)
	return nil
}

func taskDefinitionFromWorkerType(workerType WorkerType) string {
	switch workerType {
	case WorkerTypeScraper:
		return "scraper-worker-node"
	case WorkerTypeRag:
		return "rag-worker-node"
	}
	panic(fmt.Sprintf("Invalid worker type: %s", workerType))
}

func containerOverridesFrom(workerType WorkerType, params *WorkerNodeParams) ([]types.ContainerOverride, error) {
	taskDefinitionName := taskDefinitionFromWorkerType(workerType)

	switch workerType {
	case WorkerTypeScraper:
		if params.Concurrency < 1 {
			return nil, fmt.Errorf("concurrency must be at least 1")
		}
		return []types.ContainerOverride{
			{
				Name: aws.String(taskDefinitionName),
				Command: []string{
					"-worker=scraper",
					fmt.Sprintf("-concurrency=%d", params.Concurrency),
					fmt.Sprintf("-redis-addr=%s:%s", params.RedisAddress, params.RedisPort),
					fmt.Sprintf("-redis-password=%s", params.RedisPassword),
					fmt.Sprintf("-redis-db=%d", params.RedisDB),
				},
			},
		}, nil

	case WorkerTypeRag:
		return []types.ContainerOverride{
			{
				Name: aws.String(taskDefinitionName),
				Command: []string{
					"-worker=rag",
					fmt.Sprintf("-redis-addr=%s:%s", params.RedisAddress, params.RedisPort),
					fmt.Sprintf("-redis-password=%s", params.RedisPassword),
					fmt.Sprintf("-redis-db=%d", params.RedisDB),
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("invalid worker type: %s", workerType)
}

func logInfo(output *ecs.RunTaskOutput, workerType WorkerType, params *WorkerNodeParams) {
	if len(output.Tasks) > 0 {
		fmt.Printf("Successfully launched %d %s worker node(s)\n", params.NumberOfNodes, workerType)
		for i, task := range output.Tasks {
			fmt.Printf("Task %d ARN: %s\n", i+1, *task.TaskArn)
		}
	}

	if len(output.Failures) > 0 {
		fmt.Println("\nWarning: Some tasks failed to launch:")
		for _, failure := range output.Failures {
			fmt.Printf("- %s: %s\n", *failure.Arn, *failure.Reason)
		}
	}
}
