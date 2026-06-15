variable "env" {
  type        = string
  description = "Deployment environment (e.g. dev, staging, prod)"
}

variable "project" {
  type        = string
  description = "Project name"
  default     = "digimap"
}

variable "app_image" {
  type        = string
  description = "Full ECR image URI with tag (e.g. 123456789.dkr.ecr.ap-northeast-1.amazonaws.com/digimap-backend:latest)"
}

variable "app_port" {
  type        = number
  description = "Container port the application listens on"
  default     = 8080
}

variable "cpu" {
  type        = string
  description = "Task CPU units (e.g. 256, 512, 1024)"
  default     = "256"
}

variable "memory" {
  type        = string
  description = "Task memory in MiB (e.g. 512, 1024, 2048)"
  default     = "512"
}

variable "desired_count" {
  type        = number
  description = "Desired number of running tasks"
  default     = 2
}

variable "execution_role_arn" {
  type        = string
  description = "ARN of the ECS task execution role"
}

variable "task_role_arn" {
  type        = string
  description = "ARN of the ECS task role"
}

variable "private_subnet_ids" {
  type        = list(string)
  description = "Private subnet IDs for the ECS service"
}

variable "security_group_id" {
  type        = string
  description = "Security group ID for the ECS tasks"
}

variable "target_group_blue_arn" {
  type        = string
  description = "ARN of the blue target group"
}

variable "target_group_green_arn" {
  type        = string
  description = "ARN of the green target group"
}

variable "target_group_blue_name" {
  type        = string
  description = "Name of the blue target group (used by CodeDeploy)"
}

variable "target_group_green_name" {
  type        = string
  description = "Name of the green target group (used by CodeDeploy)"
}

variable "https_listener_arn" {
  type        = string
  description = "ARN of the HTTPS ALB listener"
}

variable "log_retention_days" {
  type        = number
  description = "CloudWatch log retention in days"
  default     = 30
}

variable "aws_region" {
  type        = string
  description = "AWS region"
  default     = "ap-northeast-1"
}

variable "deployment_timeout_minutes" {
  type        = number
  description = "CodeDeploy deployment timeout in minutes"
  default     = 15
}

# ─── Application environment variables ───────────────────────────────────────

variable "db_host" {
  type        = string
  description = "PostgreSQL hostname"
}

variable "db_port" {
  type        = string
  description = "PostgreSQL port"
  default     = "5432"
}

variable "db_name" {
  type        = string
  description = "PostgreSQL database name"
}

variable "db_username" {
  type        = string
  description = "PostgreSQL username"
}

variable "db_sslmode" {
  type        = string
  description = "PostgreSQL SSL mode"
  default     = "require"
}

variable "db_max_conns" {
  type        = string
  description = "Maximum database connections"
  default     = "50"
}

variable "db_min_conns" {
  type        = string
  description = "Minimum database connections"
  default     = "10"
}

variable "redis_addr" {
  type        = string
  description = "Redis address in host:port format"
}

variable "s3_bucket" {
  type        = string
  description = "Primary S3 bucket name"
}

variable "s3_assets_bucket" {
  type        = string
  description = "S3 assets bucket name"
}

variable "s3_snapshot_bucket" {
  type        = string
  description = "S3 snapshot bucket name"
}

variable "jwt_issuer" {
  type        = string
  description = "JWT issuer claim"
}

variable "jwt_access_expiry" {
  type        = string
  description = "JWT access token expiry duration"
  default     = "15m"
}

variable "jwt_refresh_expiry" {
  type        = string
  description = "JWT refresh token expiry duration"
  default     = "7d"
}

variable "log_level" {
  type        = string
  description = "Application log level"
  default     = "warn"
}

variable "log_format" {
  type        = string
  description = "Application log format"
  default     = "json"
}

# ─── Secrets Manager ARNs ────────────────────────────────────────────────────

variable "db_password_secret_arn" {
  type        = string
  description = "Secrets Manager ARN for the database password"
}

variable "jwt_secret_arn" {
  type        = string
  description = "Secrets Manager ARN for the JWT secret key"
}

variable "redis_password_secret_arn" {
  type        = string
  description = "Secrets Manager ARN for the Redis AUTH token"
}
