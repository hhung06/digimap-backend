variable "env" {
  type        = string
  description = "Deployment environment (e.g. dev, staging, prod)"
}

variable "project" {
  type        = string
  description = "Project name"
  default     = "digimap"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID where the ALB is deployed"
}

variable "public_subnet_ids" {
  type        = list(string)
  description = "Public subnet IDs for the ALB"
}

variable "security_group_id" {
  type        = string
  description = "Security group ID for the ALB"
}

variable "domain_name" {
  type        = string
  description = "API domain name (e.g. api-dev.digimap.ai)"
}

variable "zone_id" {
  type        = string
  description = "Route53 hosted zone ID"
}

variable "app_port" {
  type        = number
  description = "Application port"
  default     = 8080
}

variable "health_check_path" {
  type        = string
  description = "Health check path"
  default     = "/health"
}

variable "health_check_healthy_threshold" {
  type        = number
  description = "Number of consecutive successful health checks required"
  default     = 2
}

variable "health_check_unhealthy_threshold" {
  type        = number
  description = "Number of consecutive failed health checks required"
  default     = 3
}

variable "health_check_interval" {
  type        = number
  description = "Health check interval in seconds"
  default     = 30
}

variable "deregistration_delay" {
  type        = number
  description = "Deregistration delay in seconds for target groups"
  default     = 30
}
