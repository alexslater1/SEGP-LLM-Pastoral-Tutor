resource "aws_cloudwatch_log_group" "worker_node" {
  name              = "/ecs/worker-node"
  retention_in_days = 1 # Only keep logs for 1 day since we're doing periodic bursts
}
