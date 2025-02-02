variable "redis_password" {
  description = "Password for Redis authentication"
  type        = string
  sensitive   = true # Marks this variable as sensitive in logs
}
