# ECR Repository for Worker (used by both RAG and Scraper)
resource "aws_ecr_repository" "worker" {
  name                 = "worker-service"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  force_delete = true # Be careful with this in production
}

# ECR Repository for Queue API
resource "aws_ecr_repository" "queue_api" {
  name                 = "queue-api"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  force_delete = true # Be careful with this in production
}

# Lifecycle policy for Worker repository
resource "aws_ecr_lifecycle_policy" "worker" {
  repository = aws_ecr_repository.worker.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 5 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 5
      }
      action = {
        type = "expire"
      }
    }]
  })
}

# Lifecycle policy for Queue API repository
resource "aws_ecr_lifecycle_policy" "queue_api" {
  repository = aws_ecr_repository.queue_api.name

  policy = jsonencode({
    rules = [{
      rulePriority = 1
      description  = "Keep last 5 images"
      selection = {
        tagStatus   = "any"
        countType   = "imageCountMoreThan"
        countNumber = 5
      }
      action = {
        type = "expire"
      }
    }]
  })
}
