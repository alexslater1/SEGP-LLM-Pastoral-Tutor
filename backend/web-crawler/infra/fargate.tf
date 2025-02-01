# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "worker-node-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled" # Save costs since we don't need detailed monitoring
  }
}

# Scraper Task Definition
resource "aws_ecs_task_definition" "scraper_worker_node" {
  family                   = "scraper-worker-node"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 512  # 0.5 vCPU for web scraping
  memory                   = 1024 # 1GB RAM
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "scraper-worker-node"
      image = "${aws_ecr_repository.worker_node.repository_url}:latest"

      environment = [
        {
          name  = "WORKER_TYPE"
          value = "scraper"
        },
        {
          name  = "SCRAPE_TIMEOUT"
          value = "300" # 5 minutes in seconds
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.worker_node.name
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "scraper"
        }
      }
    }
  ])
}

# RAG Task Definition
resource "aws_ecs_task_definition" "rag_worker_node" {
  family                   = "rag-worker-node"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 1024 # 1.0 vCPU for RAG processing
  memory                   = 2048 # 2GB RAM for larger model operations
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "rag-worker-node"
      image = "${aws_ecr_repository.worker_node.repository_url}:latest"

      environment = [
        {
          name  = "WORKER_TYPE"
          value = "rag"
        },
        {
          name  = "PROCESSING_TIMEOUT"
          value = "300" # 5 minutes in seconds
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.worker_node.name
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "rag"
        }
      }
    }
  ])
}

# IAM Roles
resource "aws_iam_role" "ecs_execution_role" {
  name = "ecs-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution_role_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "ecs_task_role" {
  name = "ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

# Data source for current region
data "aws_region" "current" {}
