# Service Discovery Private DNS Namespace
resource "aws_service_discovery_private_dns_namespace" "main" {
  name        = "worker-node.local"
  description = "Private DNS namespace for worker node services"
  vpc         = aws_vpc.main.id
}
