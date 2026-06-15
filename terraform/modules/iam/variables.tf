variable "env" {
  type        = string
  description = "Deployment environment (e.g. dev, staging, prod)"
}

variable "project" {
  type        = string
  description = "Project name"
  default     = "digimap"
}

variable "secret_arns" {
  type        = list(string)
  description = "Secrets Manager ARNs the ECS task execution role needs to read"
}

variable "s3_bucket_arns" {
  type        = list(string)
  description = "S3 bucket ARNs the ECS task role needs access to"
}
