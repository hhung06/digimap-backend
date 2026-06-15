output "alb_dns_name" {
  description = "DNS name of the ALB"
  value       = aws_lb.this.dns_name
}

output "alb_zone_id" {
  description = "Hosted zone ID of the ALB (for Route53 alias records)"
  value       = aws_lb.this.zone_id
}

output "target_group_blue_arn" {
  description = "ARN of the blue target group"
  value       = aws_lb_target_group.blue.arn
}

output "target_group_green_arn" {
  description = "ARN of the green target group"
  value       = aws_lb_target_group.green.arn
}

output "target_group_blue_name" {
  description = "Name of the blue target group (used by CodeDeploy)"
  value       = aws_lb_target_group.blue.name
}

output "target_group_green_name" {
  description = "Name of the green target group (used by CodeDeploy)"
  value       = aws_lb_target_group.green.name
}

output "https_listener_arn" {
  description = "ARN of the HTTPS listener"
  value       = aws_lb_listener.https.arn
}

output "alb_arn" {
  description = "ARN of the ALB"
  value       = aws_lb.this.arn
}
