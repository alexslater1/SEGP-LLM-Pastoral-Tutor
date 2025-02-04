# Lambda function to publish Redis queue metrics
resource "aws_lambda_function" "queue_metrics" {
  filename         = "queue_metrics.zip"
  source_code_hash = filebase64sha256("queue_metrics.zip")
  function_name    = "redis-queue-metrics"
  role             = aws_iam_role.lambda_role.arn
  handler          = "dist/index.handler"
  runtime          = "nodejs22.x"
  timeout          = 60
  memory_size      = 256

  vpc_config {
    subnet_ids         = aws_subnet.public[*].id
    security_group_ids = [aws_security_group.worker_node.id]
  }

  environment {
    variables = {
      REDIS_HOST     = aws_elasticache_replication_group.queue.primary_endpoint_address
      REDIS_PORT     = "6379"
      REDIS_PASSWORD = var.redis_auth_token
    }
  }
}

# CloudWatch Event to trigger the Lambda every minute
resource "aws_cloudwatch_event_rule" "every_minute" {
  name                = "every-minute"
  description         = "Fires every minute"
  schedule_expression = "rate(1 minute)"
}

resource "aws_cloudwatch_event_target" "check_queues_every_minute" {
  rule      = aws_cloudwatch_event_rule.every_minute.name
  target_id = "CheckQueues"
  arn       = aws_lambda_function.queue_metrics.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.queue_metrics.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.every_minute.arn
}
