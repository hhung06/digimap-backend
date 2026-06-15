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
  default = "develop"
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

# ── Database ───────────────────────────────────────────────────────────────────

variable "db_instance_class" {
  type    = string
  default = "db.t4g.micro"
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
  default = "cache.t4g.micro"
}

# ── EC2 ────────────────────────────────────────────────────────────────────────

variable "instance_type" {
  type    = string
  default = "t3.small"
}

variable "key_name" {
  type        = string
  description = "EC2 key pair name for SSH. Leave empty to skip."
  default     = ""
}

variable "allowed_ssh_cidrs" {
  type        = list(string)
  description = "CIDRs allowed SSH access. Set to your IP (e.g. [\"1.2.3.4/32\"])."
  default     = []
}

variable "app_image" {
  type        = string
  description = "Full ECR image URI with tag"
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
