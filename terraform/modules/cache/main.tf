terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}

locals {
  name_prefix = "${var.project}-${var.env}"
}

resource "random_password" "redis_auth_token" {
  length           = 32
  special          = false
  override_special = ""
}

resource "aws_secretsmanager_secret" "redis_auth_token" {
  name        = "${local.name_prefix}/redis-auth-token"
  description = "Redis AUTH token for ${local.name_prefix}"

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_secretsmanager_secret_version" "redis_auth_token" {
  secret_id     = aws_secretsmanager_secret.redis_auth_token.id
  secret_string = random_password.redis_auth_token.result
}

resource "aws_elasticache_subnet_group" "this" {
  name       = local.name_prefix
  subnet_ids = var.subnet_ids

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_elasticache_replication_group" "this" {
  replication_group_id = local.name_prefix
  description          = "${var.project} ${var.env} Redis"

  node_type          = var.node_type
  num_cache_clusters = var.num_cache_clusters
  engine             = "redis"
  engine_version     = "7.1"
  parameter_group_name = "default.redis7"
  port               = 6379

  automatic_failover_enabled = var.automatic_failover_enabled

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  auth_token                 = random_password.redis_auth_token.result

  subnet_group_name  = aws_elasticache_subnet_group.this.name
  security_group_ids = var.security_group_ids

  snapshot_retention_limit = 1

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}
