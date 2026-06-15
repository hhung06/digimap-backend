variable "env" {
  type        = string
  description = "Deployment environment (e.g. dev, staging, prod)"
}

variable "project" {
  type        = string
  description = "Project name"
  default     = "digimap"
}

variable "node_type" {
  type        = string
  description = "ElastiCache node type (e.g. cache.t4g.micro)"
}

variable "num_cache_clusters" {
  type        = number
  description = "Number of cache clusters. Set to 2 with automatic_failover_enabled = true for prod."
  default     = 1
}

variable "subnet_ids" {
  type        = list(string)
  description = "Private subnet IDs for the ElastiCache subnet group"
}

variable "security_group_ids" {
  type        = list(string)
  description = "Security group IDs to attach to the replication group"
}

variable "automatic_failover_enabled" {
  type        = bool
  description = "Enable automatic failover (requires num_cache_clusters >= 2)"
  default     = false
}
