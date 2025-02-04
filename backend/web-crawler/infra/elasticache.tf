# ElastiCache subnet group
resource "aws_elasticache_subnet_group" "redis" {
  name       = "redis-cache-subnet"
  subnet_ids = aws_subnet.public[*].id
}

# ElastiCache parameter group
resource "aws_elasticache_parameter_group" "redis" {
  family = "redis7"
  name   = "redis-params"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }
}

# Replace the cluster with a replication group
resource "aws_elasticache_replication_group" "queue" {
  replication_group_id = "worker-queue"
  description          = "Worker queue redis"
  node_type            = "cache.t4g.micro"
  num_cache_clusters   = 1
  parameter_group_name = aws_elasticache_parameter_group.redis.name
  port                 = 6379
  security_group_ids   = [aws_security_group.redis.id]
  subnet_group_name    = aws_elasticache_subnet_group.redis.name

  auth_token                 = var.redis_auth_token
  transit_encryption_enabled = true
  at_rest_encryption_enabled = true
}
