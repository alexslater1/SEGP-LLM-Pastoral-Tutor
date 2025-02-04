variable "redis_auth_token" {
  description = "Auth token for Redis authentication (must be at least 16 characters)"
  type        = string
  sensitive   = true
}
