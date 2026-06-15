terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  required_version = "~> 1.9"
}

resource "aws_db_subnet_group" "this" {
  name       = "${var.project}-${var.env}"
  subnet_ids = var.subnet_ids

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_db_parameter_group" "this" {
  name   = "${var.project}-${var.env}-pg17"
  family = "postgres17"

  # Extensions are created by migrations; pgvector needs no preload entry.
  parameter {
    name  = "shared_preload_libraries"
    value = ""
  }

  parameter {
    name  = "log_connections"
    value = "1"
  }

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_db_instance" "this" {
  identifier = "${var.project}-${var.env}"

  engine         = "postgres"
  engine_version = "17"

  instance_class        = var.instance_class
  multi_az              = var.multi_az
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.this.name
  parameter_group_name   = aws_db_parameter_group.this.name
  vpc_security_group_ids = var.security_group_ids

  storage_encrypted = true
  storage_type      = "gp3"

  backup_retention_period = var.backup_retention_days
  backup_window           = "03:00-04:00"
  maintenance_window      = "Mon:04:00-Mon:05:00"

  deletion_protection      = var.deletion_protection
  skip_final_snapshot      = !var.deletion_protection
  final_snapshot_identifier = "${var.project}-${var.env}-final"

  apply_immediately           = false
  auto_minor_version_upgrade  = true
  performance_insights_enabled = true

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}
