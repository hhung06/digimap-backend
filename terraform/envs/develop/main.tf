terraform {
  required_version = ">= 1.10.0"

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

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project   = var.project
      Env       = var.env
      ManagedBy = "terraform"
    }
  }
}

data "aws_caller_identity" "current" {}

# ── Network ────────────────────────────────────────────────────────────────────

module "network" {
  source = "../../modules/network"

  env                  = var.env
  project              = var.project
  vpc_cidr             = var.vpc_cidr
  azs                  = var.azs
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  single_nat_gateway   = true
}

# ── Security Groups (for RDS + Redis) ─────────────────────────────────────────

module "security" {
  source = "../../modules/security"

  env      = var.env
  project  = var.project
  vpc_id   = module.network.vpc_id
  vpc_cidr = module.network.vpc_cidr_block
}

# ── ECR ────────────────────────────────────────────────────────────────────────

module "ecr" {
  source = "../../modules/ecr"

  env         = var.env
  project     = var.project
  image_count = 10
}

# ── Secrets Manager ────────────────────────────────────────────────────────────

module "secrets" {
  source = "../../modules/secrets"

  env        = var.env
  project    = var.project
  jwt_secret = var.jwt_secret
}

# ── Database (RDS PostgreSQL 17) ───────────────────────────────────────────────

module "database" {
  source = "../../modules/database"

  env                   = var.env
  project               = var.project
  instance_class        = var.db_instance_class
  multi_az              = false
  subnet_ids            = module.network.private_subnet_ids
  security_group_ids    = [module.security.rds_sg_id]
  db_name               = var.db_name
  db_username           = var.db_username
  db_password           = module.secrets.db_password
  allocated_storage     = 20
  max_allocated_storage = 50
  backup_retention_days = 3
  deletion_protection   = false
}

# ── Cache (ElastiCache Redis 7) ────────────────────────────────────────────────

module "cache" {
  source = "../../modules/cache"

  env                        = var.env
  project                    = var.project
  node_type                  = var.redis_node_type
  num_cache_clusters         = 1
  automatic_failover_enabled = false
  subnet_ids                 = module.network.private_subnet_ids
  security_group_ids         = [module.security.redis_sg_id]
}

# ── S3 Buckets ─────────────────────────────────────────────────────────────────

module "storage" {
  source = "../../modules/storage"

  env           = var.env
  project       = var.project
  force_destroy = true
}

# ── EC2 App Instance ───────────────────────────────────────────────────────────

module "ec2" {
  source = "../../modules/ec2-app"

  env       = var.env
  project   = var.project
  aws_region = var.aws_region

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_ids[0]

  instance_type     = var.instance_type
  key_name          = var.key_name
  allowed_ssh_cidrs = var.allowed_ssh_cidrs
  app_image         = var.app_image
  ecr_registry      = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"

  db_host     = module.database.endpoint
  db_port     = tostring(module.database.port)
  db_name     = module.database.db_name
  db_username = module.database.db_username
  db_sslmode  = "require"
  redis_addr  = "${module.cache.primary_endpoint_address}:${module.cache.port}"

  s3_bucket          = module.storage.main_bucket_name
  s3_assets_bucket   = module.storage.assets_bucket_name
  s3_snapshot_bucket = module.storage.snapshot_bucket_name
  s3_bucket_arns     = module.storage.bucket_arns

  jwt_issuer         = var.jwt_issuer
  jwt_access_expiry  = var.jwt_access_expiry
  jwt_refresh_expiry = var.jwt_refresh_expiry

  db_password_secret_arn    = module.secrets.db_password_secret_arn
  jwt_secret_arn            = module.secrets.jwt_secret_arn
  redis_password_secret_arn = module.cache.redis_auth_token_secret_arn
}

# ── Outputs ────────────────────────────────────────────────────────────────────

output "ecr_repository_url" {
  value = module.ecr.repository_url
}

output "app_public_ip" {
  description = "EC2 instance public IP — access the API at http://<ip>:8080"
  value       = module.ec2.public_ip
}

output "instance_id" {
  value = module.ec2.instance_id
}
