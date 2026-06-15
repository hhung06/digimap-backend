variable "env" {
  type = string
}

variable "project" {
  type    = string
  default = "digimap"
}

variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  type        = string
  description = "Public subnet to place the instance in"
}

variable "instance_type" {
  type    = string
  default = "t3.small"
}

variable "key_name" {
  type        = string
  description = "EC2 key pair name for SSH access. Leave empty to skip key attachment."
  default     = ""
}

variable "allowed_ssh_cidrs" {
  type        = list(string)
  description = "CIDRs allowed to SSH (port 22). Set to your IP."
  default     = []
}

variable "app_port" {
  type    = number
  default = 8080
}

# ── App image ──────────────────────────────────────────────────────────────────

variable "app_image" {
  type        = string
  description = "Full ECR image URI with tag"
}

variable "ecr_registry" {
  type        = string
  description = "ECR registry base URL (e.g. 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com)"
}

# ── App config ─────────────────────────────────────────────────────────────────

variable "db_host" { type = string }
variable "db_port" { type = string; default = "5432" }
variable "db_name" { type = string }
variable "db_username" { type = string }
variable "db_sslmode" { type = string; default = "require" }
variable "redis_addr" { type = string }

variable "s3_bucket" { type = string }
variable "s3_assets_bucket" { type = string }
variable "s3_snapshot_bucket" { type = string }

variable "jwt_issuer" { type = string }
variable "jwt_access_expiry" { type = string; default = "15m" }
variable "jwt_refresh_expiry" { type = string; default = "7d" }

# ── Secrets ────────────────────────────────────────────────────────────────────

variable "db_password_secret_arn" { type = string }
variable "jwt_secret_arn" { type = string }
variable "redis_password_secret_arn" { type = string }

variable "s3_bucket_arns" {
  type        = list(string)
  description = "ARNs of all S3 buckets the instance needs access to"
}
