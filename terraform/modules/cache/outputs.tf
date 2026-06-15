output "primary_endpoint_address" {
  description = "Redis primary endpoint address"
  value       = aws_elasticache_replication_group.this.primary_endpoint_address
}

output "port" {
  description = "Redis port"
  value       = 6379
}

output "redis_auth_token" {
  description = "Redis AUTH token (sensitive)"
  value       = random_password.redis_auth_token.result
  sensitive   = true
}

output "redis_auth_token_secret_arn" {
  description = "ARN of the Secrets Manager secret holding the Redis AUTH token"
  value       = aws_secretsmanager_secret.redis_auth_token.arn
}
