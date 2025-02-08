# Security Group for ALB
resource "aws_security_group" "alb" {
  name        = "queue-api-alb-sg"
  description = "Security group for Queue API ALB"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "queue-api-alb-sg"
  }
}

# Application Load Balancer
resource "aws_lb" "queue_api" {
  name               = "queue-api-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = [aws_subnet.public.id]

  enable_deletion_protection = false # Set to true for production

  tags = {
    Name = "queue-api-alb"
  }
}

# Target Group
resource "aws_lb_target_group" "queue_api" {
  name        = "queue-api-tg"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip" # Since we'll be using ECS Fargate

  health_check {
    enabled             = true
    healthy_threshold   = 2
    interval            = 30
    matcher             = "200"
    path                = "/ping"
    port                = "traffic-port"
    protocol            = "HTTP"
    timeout             = 5
    unhealthy_threshold = 2
  }
}

# Listener
resource "aws_lb_listener" "queue_api" {
  load_balancer_arn = aws_lb.queue_api.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.queue_api.arn
  }
}
