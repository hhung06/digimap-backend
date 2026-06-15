variable "env" {
  description = "Deployment environment (e.g. develop, staging, production)"
  type        = string
}

variable "project" {
  description = "Project name used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "jwt_secret" {
  description = "JWT signing key value stored in Secrets Manager"
  type        = string
  sensitive   = true
}
