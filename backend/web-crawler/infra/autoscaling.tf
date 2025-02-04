# Application Auto Scaling Target for Scraper Workers
resource "aws_appautoscaling_target" "scraper_workers" {
  max_capacity       = 10
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.main.name}/scraper-workers"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

# Application Auto Scaling Target for RAG Workers
resource "aws_appautoscaling_target" "rag_workers" {
  max_capacity       = 5
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.main.name}/rag-workers"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

# CloudWatch Metric Alarm for Scraper Queue
resource "aws_cloudwatch_metric_alarm" "scraper_queue_high" {
  alarm_name          = "scraper-queue-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ScraperQueueLength"
  namespace           = "CustomRedisMetrics"
  period              = "60"
  statistic           = "Average"
  threshold           = "10" # Adjust based on your needs
  alarm_description   = "This metric monitors scraper queue length"
  alarm_actions       = [aws_appautoscaling_policy.scraper_scale_up.arn]

  dimensions = {
    QueueName = "urls"
  }
}

# CloudWatch Metric Alarm for RAG Queue
resource "aws_cloudwatch_metric_alarm" "rag_queue_high" {
  alarm_name          = "rag-queue-high"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "RAGQueueLength"
  namespace           = "CustomRedisMetrics"
  period              = "60"
  statistic           = "Average"
  threshold           = "5" # Adjust based on your needs
  alarm_description   = "This metric monitors RAG queue length"
  alarm_actions       = [aws_appautoscaling_policy.rag_scale_up.arn]

  dimensions = {
    QueueName = "rag"
  }
}

# Auto Scaling Policies
resource "aws_appautoscaling_policy" "scraper_scale_up" {
  name               = "scraper-scale-up"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.scraper_workers.resource_id
  scalable_dimension = aws_appautoscaling_target.scraper_workers.scalable_dimension
  service_namespace  = aws_appautoscaling_target.scraper_workers.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = 1
    }
  }
}

resource "aws_appautoscaling_policy" "rag_scale_up" {
  name               = "rag-scale-up"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.rag_workers.resource_id
  scalable_dimension = aws_appautoscaling_target.rag_workers.scalable_dimension
  service_namespace  = aws_appautoscaling_target.rag_workers.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_lower_bound = 0
      scaling_adjustment          = 1
    }
  }
}

# Scale down policies
resource "aws_appautoscaling_policy" "scraper_scale_down" {
  name               = "scraper-scale-down"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.scraper_workers.resource_id
  scalable_dimension = aws_appautoscaling_target.scraper_workers.scalable_dimension
  service_namespace  = aws_appautoscaling_target.scraper_workers.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_upper_bound = 0
      scaling_adjustment          = -1
    }
  }
}

resource "aws_appautoscaling_policy" "rag_scale_down" {
  name               = "rag-scale-down"
  policy_type        = "StepScaling"
  resource_id        = aws_appautoscaling_target.rag_workers.resource_id
  scalable_dimension = aws_appautoscaling_target.rag_workers.scalable_dimension
  service_namespace  = aws_appautoscaling_target.rag_workers.service_namespace

  step_scaling_policy_configuration {
    adjustment_type         = "ChangeInCapacity"
    cooldown                = 60
    metric_aggregation_type = "Average"

    step_adjustment {
      metric_interval_upper_bound = 0
      scaling_adjustment          = -1
    }
  }
}

# Scale down alarms
resource "aws_cloudwatch_metric_alarm" "scraper_queue_low" {
  alarm_name          = "scraper-queue-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "ScraperQueueLength"
  namespace           = "CustomRedisMetrics"
  period              = "60"
  statistic           = "Average"
  threshold           = "5" # Adjust based on your needs
  alarm_description   = "This metric monitors scraper queue length for scale down"
  alarm_actions       = [aws_appautoscaling_policy.scraper_scale_down.arn]

  dimensions = {
    QueueName = "urls"
  }
}

resource "aws_cloudwatch_metric_alarm" "rag_queue_low" {
  alarm_name          = "rag-queue-low"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "2"
  metric_name         = "RAGQueueLength"
  namespace           = "CustomRedisMetrics"
  period              = "60"
  statistic           = "Average"
  threshold           = "2" # Adjust based on your needs
  alarm_description   = "This metric monitors RAG queue length for scale down"
  alarm_actions       = [aws_appautoscaling_policy.rag_scale_down.arn]

  dimensions = {
    QueueName = "rag"
  }
}
