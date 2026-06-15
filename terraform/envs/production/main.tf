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

# ── Network ────────────────────────────────────────────────────────────────────

module "network" {
  source = "../../modules/network"

  env                  = var.env
  project              = var.project
  vpc_cidr             = var.vpc_cidr
  azs                  = var.azs
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  single_nat_gateway   = false # one NAT per AZ for production HA
}

# ── Security Groups ────────────────────────────────────────────────────────────

module "security" {
  source = "../../modules/security"

  env      = var.env
  project  = var.project
  vpc_id   = module.network.vpc_id
  vpc_cidr = module.network.vpc_cidr_block
}

# ── ECR (shared with develop — one repo, env-tagged images) ───────────────────

module "ecr" {
  source = "../../modules/ecr"

  env         = var.env
  project     = var.project
  image_count = 30
}

# ── Secrets Manager ────────────────────────────────────────────────────────────

module "secrets" {
  source = "../../modules/secrets"

  env        = var.env
  project    = var.project
  jwt_secret = var.jwt_secret
}

# ── Database (RDS PostgreSQL 17, Multi-AZ) ─────────────────────────────────────

module "database" {
  source = "../../modules/database"

  env                   = var.env
  project               = var.project
  instance_class        = var.db_instance_class
  multi_az              = true
  subnet_ids            = module.network.private_subnet_ids
  security_group_ids    = [module.security.rds_sg_id]
  db_name               = var.db_name
  db_username           = var.db_username
  db_password           = module.secrets.db_password
  allocated_storage     = 50
  max_allocated_storage = 200
  backup_retention_days = 7
  deletion_protection   = true
}

# ── Cache (ElastiCache Redis 7, 2-node with failover) ─────────────────────────

module "cache" {
  source = "../../modules/cache"

  env                        = var.env
  project                    = var.project
  node_type                  = var.redis_node_type
  num_cache_clusters         = 2
  automatic_failover_enabled = true
  subnet_ids                 = module.network.private_subnet_ids
  security_group_ids         = [module.security.redis_sg_id]
}

# ── S3 Buckets ─────────────────────────────────────────────────────────────────

module "storage" {
  source = "../../modules/storage"

  env           = var.env
  project       = var.project
  force_destroy = false # protect production data
}

# ── IAM Roles ─────────────────────────────────────────────────────────────────

module "iam" {
  source = "../../modules/iam"

  env     = var.env
  project = var.project

  secret_arns = [
    module.secrets.db_password_secret_arn,
    module.secrets.jwt_secret_arn,
    module.cache.redis_auth_token_secret_arn,
  ]

  s3_bucket_arns = module.storage.bucket_arns
}

# ── ALB + ACM ─────────────────────────────────────────────────────────────────

module "alb" {
  source = "../../modules/alb"

  env               = var.env
  project           = var.project
  vpc_id            = module.network.vpc_id
  public_subnet_ids = module.network.public_subnet_ids
  security_group_id = module.security.alb_sg_id
  domain_name       = var.domain_name
  zone_id           = var.route53_zone_id

  # Tighten health check thresholds for production
  health_check_healthy_threshold   = 2
  health_check_unhealthy_threshold = 2
  deregistration_delay             = 60
}

# ── ECS Fargate + CodeDeploy ──────────────────────────────────────────────────

module "ecs" {
  source = "../../modules/ecs-service"

  env     = var.env
  project = var.project

  # Image
  app_image = var.app_image

  # Sizing
  cpu           = var.app_cpu
  memory        = var.app_memory
  desired_count = var.app_desired_count

  # IAM
  execution_role_arn = module.iam.execution_role_arn
  task_role_arn      = module.iam.task_role_arn

  # Networking
  private_subnet_ids = module.network.private_subnet_ids
  security_group_id  = module.security.app_sg_id

  # ALB
  target_group_blue_arn   = module.alb.target_group_blue_arn
  target_group_green_arn  = module.alb.target_group_green_arn
  target_group_blue_name  = module.alb.target_group_blue_name
  target_group_green_name = module.alb.target_group_green_name
  https_listener_arn      = module.alb.https_listener_arn

  # App config
  aws_region   = var.aws_region
  db_host      = module.database.endpoint
  db_port      = tostring(module.database.port)
  db_name      = module.database.db_name
  db_username  = module.database.db_username
  db_sslmode   = "require"
  db_max_conns = "50"
  db_min_conns = "10"

  redis_addr = "${module.cache.primary_endpoint_address}:${module.cache.port}"

  s3_bucket          = module.storage.main_bucket_name
  s3_assets_bucket   = module.storage.assets_bucket_name
  s3_snapshot_bucket = module.storage.snapshot_bucket_name

  jwt_issuer         = var.jwt_issuer
  jwt_access_expiry  = var.jwt_access_expiry
  jwt_refresh_expiry = var.jwt_refresh_expiry

  log_level  = "warn"
  log_format = "json"

  log_retention_days         = 90
  deployment_timeout_minutes = 20

  # Secrets (Secrets Manager ARNs)
  db_password_secret_arn    = module.secrets.db_password_secret_arn
  jwt_secret_arn            = module.secrets.jwt_secret_arn
  redis_password_secret_arn = module.cache.redis_auth_token_secret_arn
}

# ── Outputs ────────────────────────────────────────────────────────────────────

output "ecr_repository_url" {
  value = module.ecr.repository_url
}

output "alb_dns_name" {
  value = module.alb.alb_dns_name
}

output "codedeploy_app_name" {
  value = module.ecs.codedeploy_app_name
}

output "codedeploy_deployment_group_name" {
  value = module.ecs.codedeploy_deployment_group_name
}

output "ecs_cluster_name" {
  value = module.ecs.cluster_name
}

output "ecs_service_name" {
  value = module.ecs.service_name
}
