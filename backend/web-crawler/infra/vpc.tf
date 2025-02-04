# VPC
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "worker-node-vpc"
  }
}

# Public Subnets (2 for high availability)
resource "aws_subnet" "public" {
  count             = 2
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.${count.index + 1}.0/24"
  availability_zone = data.aws_availability_zones.available.names[count.index]

  map_public_ip_on_launch = true

  tags = {
    Name = "worker-node-public-${count.index + 1}"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "worker-node-igw"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "worker-node-public-rt"
  }
}

# Route Table Association
resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# Security Group
resource "aws_security_group" "worker_node" {
  name        = "worker-node-sg"
  description = "Security group for worker node"
  vpc_id      = aws_vpc.main.id

  # No inbound rules needed since we're only making outbound requests

  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Allow HTTPS to anywhere instead of using prefix list
  }

  egress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  tags = {
    Name = "worker-node-sg"
  }
}

# Security Group for Redis
resource "aws_security_group" "redis" {
  name        = "redis-sg"
  description = "Security group for Redis"
  vpc_id      = aws_vpc.main.id

  # Allow inbound traffic on Redis port from worker nodes and Lambda
  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.worker_node.id]
  }

  # No outbound rules needed for Redis
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "redis-sg"
  }
}

# Data source for availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

# VPC Endpoint for CloudWatch
resource "aws_vpc_endpoint" "cloudwatch" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.monitoring"
  vpc_endpoint_type = "Interface"

  subnet_ids         = aws_subnet.public[*].id
  security_group_ids = [aws_security_group.worker_node.id]

  private_dns_enabled = true
}

# VPC Endpoint for CloudWatch Logs (needed for Lambda logging)
resource "aws_vpc_endpoint" "cloudwatch_logs" {
  vpc_id            = aws_vpc.main.id
  service_name      = "com.amazonaws.${data.aws_region.current.name}.logs"
  vpc_endpoint_type = "Interface"

  subnet_ids         = aws_subnet.public[*].id
  security_group_ids = [aws_security_group.worker_node.id]

  private_dns_enabled = true
}
