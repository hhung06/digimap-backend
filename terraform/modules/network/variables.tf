variable "env" {
  description = "Environment name (e.g. develop, production)"
  type        = string
}

variable "project" {
  description = "Project name, used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
}

variable "azs" {
  description = "List of exactly 2 Availability Zone names to deploy subnets into"
  type        = list(string)

  validation {
    condition     = length(var.azs) == 2
    error_message = "Exactly 2 Availability Zones must be provided."
  }
}

variable "public_subnet_cidrs" {
  description = "List of 2 CIDR blocks for the public subnets (one per AZ)"
  type        = list(string)

  validation {
    condition     = length(var.public_subnet_cidrs) == 2
    error_message = "Exactly 2 public subnet CIDRs must be provided."
  }
}

variable "private_subnet_cidrs" {
  description = "List of 2 CIDR blocks for the private subnets (one per AZ)"
  type        = list(string)

  validation {
    condition     = length(var.private_subnet_cidrs) == 2
    error_message = "Exactly 2 private subnet CIDRs must be provided."
  }
}

variable "single_nat_gateway" {
  description = "When true, a single NAT gateway is shared across all private subnets. When false, one NAT gateway is created per AZ."
  type        = bool
  default     = true
}
