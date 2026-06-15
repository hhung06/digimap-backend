variable "env" {
  description = "Environment name (e.g. develop, production)"
  type        = string
}

variable "project" {
  description = "Project name, used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "vpc_id" {
  description = "ID of the VPC in which to create the security groups"
  type        = string
}

variable "vpc_cidr" {
  description = "CIDR block of the VPC (used for internal-only rules if needed)"
  type        = string
}
