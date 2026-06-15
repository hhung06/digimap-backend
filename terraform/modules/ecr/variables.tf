variable "env" {
  description = "Deployment environment (e.g. develop, staging, production)"
  type        = string
}

variable "project" {
  description = "Project name used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "image_count" {
  description = "Maximum number of tagged images to retain in the ECR lifecycle policy"
  type        = number
  default     = 20
}
