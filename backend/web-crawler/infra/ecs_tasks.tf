# Task Definition for Queue API
resource "aws_ecs_task_definition" "queue_api" {
  family                   = "queue-api"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256" # 0.25 vCPU
  memory                   = "512" # 512 MB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "queue-api"
      image     = "${aws_ecr_repository.queue_api.repository_url}:latest"
      essential = true

      portMappings = [
        {
          containerPort = 80
          hostPort      = 80
          protocol      = "tcp"
        }
      ]

      secrets = [
        {
          name      = "REDIS_PASSWORD"
          valueFrom = aws_secretsmanager_secret.redis_password.arn
        }
      ]

      environment = [
        {
          name  = "REDIS_ADDRESS"
          value = "${aws_elasticache_replication_group.queue.primary_endpoint_address}:6379"
        },
        {
          name  = "REDIS_DB"
          value = "1"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/queue-api"
          "awslogs-region"        = "eu-west-2"
          "awslogs-stream-prefix" = "ecs"
          "awslogs-create-group"  = "true"
        }
      }
    }
  ])
}

# Task Definition for RAG Worker
resource "aws_ecs_task_definition" "rag_worker" {
  family                   = "rag-worker"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "1024" # 1 vCPU
  memory                   = "2048" # 2 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "rag-worker"
      image     = "${aws_ecr_repository.worker.repository_url}:latest"
      essential = true

      secrets = [
        {
          name      = "REDIS_PASSWORD"
          valueFrom = aws_secretsmanager_secret.redis_password.arn
        }
      ]

      environment = [
        {
          name  = "REDIS_ADDR"
          value = "${aws_elasticache_replication_group.queue.primary_endpoint_address}:6379"
        },
        {
          name  = "REDIS_DB"
          value = "1"
        },
        {
          name  = "WORKER_TYPE"
          value = "rag"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/rag-worker"
          "awslogs-region"        = "eu-west-2"
          "awslogs-stream-prefix" = "ecs"
          "awslogs-create-group"  = "true"
        }
      }
    }
  ])
}

# Task Definition for Scraper Worker
resource "aws_ecs_task_definition" "scraper_worker" {
  family                   = "scraper-worker"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"  # 0.5 vCPU
  memory                   = "1024" # 1 GB
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name      = "scraper-worker"
      image     = "${aws_ecr_repository.worker.repository_url}:latest"
      essential = true

      secrets = [
        {
          name      = "REDIS_PASSWORD"
          valueFrom = aws_secretsmanager_secret.redis_password.arn
        }
      ]

      environment = [
        {
          name  = "REDIS_ADDR"
          value = "${aws_elasticache_replication_group.queue.primary_endpoint_address}:6379"
        },
        {
          name  = "REDIS_DB"
          value = "1"
        },
        {
          name  = "WORKER_TYPE"
          value = "scraper"
        },
        {
          name  = "CONCURRENCY"
          value = "10"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/scraper-worker"
          "awslogs-region"        = "eu-west-2"
          "awslogs-stream-prefix" = "ecs"
          "awslogs-create-group"  = "true"
        }
      }
    }
  ])
}
