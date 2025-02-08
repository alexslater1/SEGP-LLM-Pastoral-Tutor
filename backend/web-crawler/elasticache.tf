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

# Elasticache Redis cluster
resource "aws_elasticache_cluster" "main" {
  cluster_id           = "redis-queue"
  engine               = "redis"
  node_type            = "cache.t4g.micro" # Smallest instance for dev, adjust for production
  num_cache_nodes      = 1
  parameter_group_name = aws_elasticache_parameter_group.main.name
  port                 = 6379
  security_group_ids   = [aws_security_group.elasticache.id]
  subnet_group_name    = aws_elasticache_subnet_group.main.name

  # Recommended settings
  apply_immediately  = true
  az_mode            = "single-az" # Use multi-az for production
  maintenance_window = "sun:05:00-sun:06:00"

  tags = {
    Name = "redis-queue"
  }
}

# Output the Redis endpoint for reference
output "redis_endpoint" {
  value = aws_elasticache_cluster.main.cache_nodes[0].address
}
