# ECR Repository URLs
output "queue_api_repository_url" {
  description = "The URL of the Queue API ECR repository"
  value       = aws_ecr_repository.queue_api.repository_url
}

output "worker_repository_url" {
  description = "The URL of the Worker ECR repository"
  value       = aws_ecr_repository.worker.repository_url
}

# Redis
output "redis_endpoint" {
  description = "The endpoint of the Redis replication group"
  value       = aws_elasticache_replication_group.queue.primary_endpoint_address
}

output "redis_port" {
  description = "The port of the Redis replication group"
  value       = "6379"
}

# CloudWatch Log Groups
output "queue_api_log_group" {
  description = "The CloudWatch log group for the Queue API"
  value       = "/ecs/queue-api"
}

output "rag_worker_log_group" {
  description = "The CloudWatch log group for the RAG worker"
  value       = "/ecs/rag-worker"
}

output "scraper_worker_log_group" {
  description = "The CloudWatch log group for the Scraper worker"
  value       = "/ecs/scraper-worker"
}

# VPC
output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.main.id
}

output "private_subnet_id" {
  description = "The ID of the private subnet"
  value       = aws_subnet.private.id
}

output "public_subnet_id" {
  description = "The ID of the public subnet"
  value       = aws_subnet.public.id
}

# ECR Push Commands
output "ecr_push_commands" {
  description = "Commands to push images to ECR"
  value = {
    login     = "aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin ${aws_ecr_repository.worker.registry_id}.dkr.ecr.eu-west-2.amazonaws.com"
    queue_api = "docker tag queue-api:latest ${aws_ecr_repository.queue_api.repository_url}:latest && docker push ${aws_ecr_repository.queue_api.repository_url}:latest"
    worker    = "docker tag worker:latest ${aws_ecr_repository.worker.repository_url}:latest && docker push ${aws_ecr_repository.worker.repository_url}:latest"
  }
}
