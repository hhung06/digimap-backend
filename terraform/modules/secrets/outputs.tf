output "db_password_secret_arn" {
  description = "ARN of the Secrets Manager secret holding the DB password"
  value       = aws_secretsmanager_secret.db_password.arn
}

output "jwt_secret_arn" {
  description = "ARN of the Secrets Manager secret holding the JWT signing key"
  value       = aws_secretsmanager_secret.jwt_secret.arn
}

output "db_password" {
  description = "Generated DB master password (sensitive)"
  value       = random_password.db_password.result
  sensitive   = true
}

output "read_secrets_policy_arn" {
  description = "ARN of the IAM policy that allows reading both secrets"
  value       = aws_iam_policy.read_secrets.arn
}
