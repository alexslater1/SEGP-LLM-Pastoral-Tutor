resource "aws_ecr_repository" "worker_node" {
  name                 = "worker-node"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

# Optional: Add a lifecycle policy to manage image versions
resource "aws_ecr_lifecycle_policy" "worker_node_policy" {
  repository = aws_ecr_repository.worker_node.name

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
