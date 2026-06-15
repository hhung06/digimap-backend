variable "env" {
  description = "Deployment environment (e.g. develop, staging, production)"
  type        = string
}

variable "project" {
  description = "Project name used as a prefix for resource names"
  type        = string
  default     = "digimap"
}

variable "instance_class" {
  description = "RDS instance class (e.g. db.t4g.micro, db.t4g.small)"
  type        = string
}

variable "multi_az" {
  description = "Enable Multi-AZ deployment for high availability"
  type        = bool
  default     = false
}

variable "subnet_ids" {
  description = "List of private subnet IDs for the DB subnet group"
  type        = list(string)
}

variable "security_group_ids" {
  description = "List of VPC security group IDs to attach to the RDS instance"
  type        = list(string)
}

variable "db_name" {
  description = "Name of the initial database to create"
  type        = string
  default     = "digimap"
}

variable "db_username" {
  description = "Master username for the RDS instance"
  type        = string
  default     = "digimap"
}

variable "db_password" {
  description = "Master password for the RDS instance"
  type        = string
  sensitive   = true
}

variable "allocated_storage" {
  description = "Initial allocated storage in GiB"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "Maximum storage autoscaling ceiling in GiB"
  type        = number
  default     = 100
}

variable "backup_retention_days" {
  description = "Number of days to retain automated backups"
  type        = number
  default     = 7
}

variable "deletion_protection" {
  description = "Enable deletion protection; also controls whether a final snapshot is taken"
  type        = bool
  default     = false
}
