# Coordinator Task Definition
resource "aws_ecs_task_definition" "coordinator" {
  family                   = "coordinator"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256 # 0.25 vCPU
  memory                   = 512 # 0.5GB RAM
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([
    {
      name  = "coordinator"
      image = "${aws_ecr_repository.coordinator.repository_url}:latest"

      portMappings = [
        {
          containerPort = 8080
          protocol      = "tcp"
        }
      ]

      environment = [
        {
          name  = "REDIS_ADDRESS"
          value = "${aws_elasticache_replication_group.queue.primary_endpoint_address}:6379"
        },
        {
          name  = "REDIS_DB"
          value = "0"
        },
        {
          name  = "REDIS_PASSWORD"
          value = var.redis_auth_token
        }
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.coordinator.name
          "awslogs-region"        = data.aws_region.current.name
          "awslogs-stream-prefix" = "coordinator"
        }
      }
    }
  ])
}

# ECR Repository for coordinator
resource "aws_ecr_repository" "coordinator" {
  name                 = "coordinator"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "coordinator" {
  name              = "/ecs/coordinator"
  retention_in_days = 14
}

# Application Load Balancer
resource "aws_lb" "coordinator" {
  name               = "coordinator-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = aws_subnet.public[*].id
}

# ALB Target Group
resource "aws_lb_target_group" "coordinator" {
  name        = "coordinator-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    path                = "/health" # Make sure to implement this endpoint
    healthy_threshold   = 2
    unhealthy_threshold = 10
  }
}

# ALB Listener
resource "aws_lb_listener" "coordinator" {
  load_balancer_arn = aws_lb.coordinator.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.coordinator.arn
  }
}

# Security Group for ALB
resource "aws_security_group" "alb" {
  name        = "coordinator-alb-sg"
  description = "Security group for coordinator ALB"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Security Group for Coordinator Service
resource "aws_security_group" "coordinator" {
  name        = "coordinator-service-sg"
  description = "Security group for coordinator service"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port       = 8080
    to_port         = 8080
    protocol        = "tcp"
    security_groups = [aws_security_group.alb.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ECS Service
resource "aws_ecs_service" "coordinator" {
  name            = "coordinator"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.coordinator.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  # Add explicit depends_on
  depends_on = [
    aws_cloudwatch_log_group.coordinator,
    aws_lb_listener.coordinator,
    aws_vpc_endpoint.cloudwatch_logs,
    aws_vpc_endpoint.ecr_api,
    aws_vpc_endpoint.ecr_dkr
  ]

  network_configuration {
    subnets          = aws_subnet.public[*].id
    security_groups  = [aws_security_group.coordinator.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.coordinator.arn
    container_name   = "coordinator"
    container_port   = 8080
  }
}

# Store Redis auth token in Secrets Manager
resource "aws_secretsmanager_secret" "redis_auth_token" {
  name = "redis-auth-token"
}

resource "aws_secretsmanager_secret_version" "redis_auth_token" {
  secret_id     = aws_secretsmanager_secret.redis_auth_token.id
  secret_string = var.redis_auth_token
}

# Update the ECS execution role policy
resource "aws_iam_role_policy" "ecs_task_execution_role_policy" {
  name = "ecs-task-execution-role-policy"
  role = aws_iam_role.ecs_execution_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage"
        ]
        Resource = "*"
      }
    ]
  })
}
