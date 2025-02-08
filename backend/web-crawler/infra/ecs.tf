# RAG Cluster
resource "aws_ecs_cluster" "rag" {
  name = "rag-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "rag-cluster"
  }
}

# Scraper Cluster
resource "aws_ecs_cluster" "scraper" {
  name = "scraper-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "scraper-cluster"
  }
}

# Queue API Cluster
resource "aws_ecs_cluster" "queue_api" {
  name = "queue-api-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Name = "queue-api-cluster"
  }
}

# Capacity Providers for RAG Cluster
resource "aws_ecs_cluster_capacity_providers" "rag" {
  cluster_name = aws_ecs_cluster.rag.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

# Capacity Providers for Scraper Cluster
resource "aws_ecs_cluster_capacity_providers" "scraper" {
  cluster_name = aws_ecs_cluster.scraper.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}

# Capacity Providers for Queue API Cluster
resource "aws_ecs_cluster_capacity_providers" "queue_api" {
  cluster_name = aws_ecs_cluster.queue_api.name

  capacity_providers = ["FARGATE"]

  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE"
  }
}
