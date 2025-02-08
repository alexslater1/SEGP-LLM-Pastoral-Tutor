# Security group for Elasticache
resource "aws_security_group" "elasticache" {
  name        = "elasticache-sg"
  description = "Security group for Elasticache Redis"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "Redis from VPC"
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.main.cidr_block] # Allow access from within VPC
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "elasticache-sg"
  }
}

# Subnet group for Elasticache
resource "aws_elasticache_subnet_group" "main" {
  name       = "elasticache-subnet-group"
  subnet_ids = [aws_subnet.private.id]
}

# Parameter group for Redis
resource "aws_elasticache_parameter_group" "main" {
  family = "redis7"
  name   = "redis-params"

  parameter {
    name  = "maxmemory-policy"
    value = "allkeys-lru"
  }

}

resource "aws_elasticache_replication_group" "queue" {
  replication_group_id = "redis-queue"
  description          = "Worker queue redis"
  node_type            = "cache.t4g.micro"
  num_cache_clusters   = 1
  parameter_group_name = aws_elasticache_parameter_group.main.name
  port                 = 6379
  security_group_ids   = [aws_security_group.elasticache.id]
  subnet_group_name    = aws_elasticache_subnet_group.main.name

  auth_token                 = var.redis_password
  transit_encryption_enabled = true
  at_rest_encryption_enabled = true
}

