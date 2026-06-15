variable "env" {
  description = "Deployment environment (e.g. develop, staging, production)"
  type        = string
}

variable "project" {
  description = "Project name used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "force_destroy" {
  description = "Allow bucket deletion even when objects are present (set true in develop)"
  type        = bool
  default     = false
}
