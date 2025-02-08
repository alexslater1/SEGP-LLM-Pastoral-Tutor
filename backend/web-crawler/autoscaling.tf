# Lambda function for queue metrics
resource "aws_lambda_function" "queue_metrics" {
  filename      = "queue_metrics.zip"
  function_name = "queue-metrics"
  role          = aws_iam_role.lambda_role.arn
  handler       = "index.handler"
  runtime       = "nodejs22.x"
  timeout       = 30

  environment {
    variables = {
      REDIS_HOST     = aws_elasticache_cluster.main.cache_nodes[0].address
      REDIS_PORT     = "6379"
      REDIS_PASSWORD = var.redis_password
    }
  }

  vpc_config {
    subnet_ids         = [aws_subnet.private.id]
    security_group_ids = [aws_security_group.lambda.id]
  }
}

# Security group for Lambda
resource "aws_security_group" "lambda" {
  name        = "lambda-sg"
  description = "Security group for Lambda function"
  vpc_id      = aws_vpc.main.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Allow Lambda to access Redis
resource "aws_security_group_rule" "lambda_to_redis" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.lambda.id
  security_group_id        = aws_security_group.elasticache.id
}

# CloudWatch Event Rule to trigger Lambda every minute
resource "aws_cloudwatch_event_rule" "queue_metrics" {
  name                = "queue-metrics-rule"
  description         = "Trigger queue metrics Lambda function"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "queue_metrics" {
  rule      = aws_cloudwatch_event_rule.queue_metrics.name
  target_id = "QueueMetricsLambda"
  arn       = aws_lambda_function.queue_metrics.arn
}

resource "aws_lambda_permission" "allow_eventbridge" {
  statement_id  = "AllowEventBridgeInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.queue_metrics.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.queue_metrics.arn
}

# Auto Scaling for RAG Service
resource "aws_appautoscaling_target" "rag" {
  max_capacity       = 5
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.rag.name}/${aws_ecs_service.rag.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "rag_queue" {
  name               = "rag-queue-policy"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.rag.resource_id
  scalable_dimension = aws_appautoscaling_target.rag.scalable_dimension
  service_namespace  = aws_appautoscaling_target.rag.service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = 10.0 # Scale when there are more than 10 items in queue
    scale_in_cooldown  = 300  # 5 minutes
    scale_out_cooldown = 60   # 1 minute

    customized_metric_specification {
      metric_name = "RAGQueueLength"
      namespace   = "CustomRedisMetrics"
      statistic   = "Average"
      unit        = "Count"

      dimensions {
        name  = "QueueName"
        value = "rag"
      }
    }
  }
}

# Auto Scaling for Scraper Service
resource "aws_appautoscaling_target" "scraper" {
  max_capacity       = 5
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.scraper.name}/${aws_ecs_service.scraper.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "scraper_queue" {
  name               = "scraper-queue-policy"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.scraper.resource_id
  scalable_dimension = aws_appautoscaling_target.scraper.scalable_dimension
  service_namespace  = aws_appautoscaling_target.scraper.service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = 20.0 # Scale when there are more than 20 items in queue
    scale_in_cooldown  = 300  # 5 minutes
    scale_out_cooldown = 60   # 1 minute

    customized_metric_specification {
      metric_name = "ScraperQueueLength"
      namespace   = "CustomRedisMetrics"
      statistic   = "Average"
      unit        = "Count"

      dimensions {
        name  = "QueueName"
        value = "urls"
      }
    }
  }
}

# Auto Scaling for Queue API based on CPU
resource "aws_appautoscaling_target" "queue_api" {
  max_capacity       = 3
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.queue_api.name}/${aws_ecs_service.queue_api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "queue_api_cpu" {
  name               = "queue-api-cpu"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.queue_api.resource_id
  scalable_dimension = aws_appautoscaling_target.queue_api.scalable_dimension
  service_namespace  = aws_appautoscaling_target.queue_api.service_namespace

  target_tracking_scaling_policy_configuration {
    target_value       = 70.0
    scale_in_cooldown  = 300
    scale_out_cooldown = 60

    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
  }
}
