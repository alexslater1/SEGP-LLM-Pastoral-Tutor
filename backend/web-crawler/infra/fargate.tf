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

# Redis Task Definition
resource "aws_ecs_task_definition" "redis" {
  family                   = "redis"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256 # 0.25 vCPU is sufficient for Redis
  memory                   = 512 # 512MB RAM
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "redis"
      image = "redis:7-alpine" # Using official Redis image, alpine for smaller size

      command = ["redis-server", "--requirepass", "${var.redis_password}"]

      portMappings = [
        {
          containerPort = 6379
          hostPort      = 6379
          protocol      = "tcp"
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.worker_node.name
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "redis"
        }
      }
    }
  ])
}

# Redis Service (keeps one instance running)
resource "aws_ecs_service" "redis" {
  name            = "redis"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.redis.arn
  desired_count   = 1 # Keep one Redis instance running
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.redis.id]
    assign_public_ip = true
  }
}

# Redis Security Group
resource "aws_security_group" "redis" {
  name        = "redis-sg"
  description = "Security group for Redis"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.worker_node.id] # Allow access from worker nodes
  }

  # Add this new ingress rule for testing from your local machine
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # WARNING: This allows access from anywhere. For testing only!
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "redis-sg"
  }
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
