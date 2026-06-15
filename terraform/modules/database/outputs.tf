output "endpoint" {
  description = "RDS instance hostname (no port)"
  value       = aws_db_instance.this.address
}

output "port" {
  description = "RDS instance port"
  value       = aws_db_instance.this.port
}

output "db_name" {
  description = "Name of the initial database"
  value       = aws_db_instance.this.db_name
}

output "db_username" {
  description = "Master username for the RDS instance"
  value       = aws_db_instance.this.username
}

output "instance_id" {
  description = "RDS instance identifier"
  value       = aws_db_instance.this.identifier
}
