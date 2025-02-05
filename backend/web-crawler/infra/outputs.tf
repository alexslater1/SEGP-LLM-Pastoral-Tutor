output "redis_endpoint" {
  value = aws_elasticache_replication_group.queue.primary_endpoint_address
}

output "coordinator_endpoint" {
  value = "http://${aws_lb.coordinator.dns_name}"
}
