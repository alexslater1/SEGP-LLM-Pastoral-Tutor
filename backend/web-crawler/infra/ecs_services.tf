# Security Groups for ECS Services
resource "aws_security_group" "queue_api" {
  name        = "queue-api-sg"
  description = "Security group for Queue API ECS Service"
  vpc_id      = aws_vpc.main.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Allow direct access from internet
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "queue-api-sg"
  }
}

resource "aws_security_group" "workers" {
  name        = "workers-sg"
  description = "Security group for RAG and Scraper workers"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "workers-sg"
  }
}

# Queue API Service
resource "aws_ecs_service" "queue_api" {
  name            = "queue-api"
  cluster         = aws_ecs_cluster.queue_api.id
  task_definition = aws_ecs_task_definition.queue_api.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.public.id]
    security_groups  = [aws_security_group.queue_api.id]
    assign_public_ip = true
  }
}

# RAG Worker Service
resource "aws_ecs_service" "rag" {
  name            = "rag-worker"
  cluster         = aws_ecs_cluster.rag.id
  task_definition = aws_ecs_task_definition.rag_worker.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.private.id]
    security_groups  = [aws_security_group.workers.id]
    assign_public_ip = false
  }
}

# Scraper Worker Service
resource "aws_ecs_service" "scraper" {
  name            = "scraper-worker"
  cluster         = aws_ecs_cluster.scraper.id
  task_definition = aws_ecs_task_definition.scraper_worker.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = [aws_subnet.private.id]
    security_groups  = [aws_security_group.workers.id]
    assign_public_ip = false
  }
}

# Allow worker security group to access Redis
resource "aws_security_group_rule" "workers_to_redis" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.workers.id
  security_group_id        = aws_security_group.elasticache.id
}

# Allow queue API security group to access Redis
resource "aws_security_group_rule" "queue_api_to_redis" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.queue_api.id
  security_group_id        = aws_security_group.elasticache.id
}
