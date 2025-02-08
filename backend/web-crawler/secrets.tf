resource "aws_secretsmanager_secret" "redis_password" {
  name = "redis-password"
}

resource "aws_secretsmanager_secret_version" "redis_password" {
  secret_id     = aws_secretsmanager_secret.redis_password.id
  secret_string = var.redis_password
}

# We'll need to manually set the secret value in AWS Secrets Manager
# or use the AWS CLI/Console to set it after creation 
