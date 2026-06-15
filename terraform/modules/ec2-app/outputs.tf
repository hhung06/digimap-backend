output "public_ip" {
  description = "Elastic IP of the app instance"
  value       = aws_eip.this.public_ip
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.this.id
}

output "security_group_id" {
  description = "Security group ID of the EC2 instance"
  value       = aws_security_group.ec2.id
}
