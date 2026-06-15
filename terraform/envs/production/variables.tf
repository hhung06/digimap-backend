variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "project" {
  type    = string
  default = "digimap"
}

variable "env" {
  type    = string
  default = "production"
}

# ── Network ────────────────────────────────────────────────────────────────────

variable "vpc_cidr" {
  type = string
}

variable "azs" {
  type = list(string)
}

variable "public_subnet_cidrs" {
  type = list(string)
}

variable "private_subnet_cidrs" {
  type = list(string)
}

# ── DNS ────────────────────────────────────────────────────────────────────────

variable "domain_name" {
  type        = string
  description = "API domain for the production environment (e.g. api.digimap.ai)"
}

variable "route53_zone_id" {
  type        = string
  description = "Route53 hosted zone ID"
}

# ── Database ───────────────────────────────────────────────────────────────────

variable "db_instance_class" {
  type    = string
  default = "db.t4g.medium"
}

variable "db_name" {
  type    = string
  default = "digimap"
}

variable "db_username" {
  type    = string
  default = "digimap"
}

# ── Cache ──────────────────────────────────────────────────────────────────────

variable "redis_node_type" {
  type    = string
  default = "cache.t4g.small"
}

# ── ECS ────────────────────────────────────────────────────────────────────────

variable "app_image" {
  type        = string
  description = "Full ECR image URI with tag — update per deployment"
}

variable "app_cpu" {
  type    = string
  default = "512"
}

variable "app_memory" {
  type    = string
  default = "1024"
}

variable "app_desired_count" {
  type    = number
  default = 2
}

# ── App config ─────────────────────────────────────────────────────────────────

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "jwt_issuer" {
  type    = string
  default = "digimap"
}

variable "jwt_access_expiry" {
  type    = string
  default = "15m"
}

variable "jwt_refresh_expiry" {
  type    = string
  default = "7d"
}
