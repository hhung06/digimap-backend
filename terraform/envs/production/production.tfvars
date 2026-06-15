# ── Core ───────────────────────────────────────────────────────────────────────
aws_region = "ap-northeast-1"
project    = "digimap"
env        = "production"

# ── Network ────────────────────────────────────────────────────────────────────
vpc_cidr             = "10.1.0.0/16"
azs                  = ["ap-northeast-1a", "ap-northeast-1c"]
public_subnet_cidrs  = ["10.1.1.0/24", "10.1.2.0/24"]
private_subnet_cidrs = ["10.1.11.0/24", "10.1.12.0/24"]

# ── DNS (update with your Route53 zone) ───────────────────────────────────────
domain_name     = "api.digimap.ai"
route53_zone_id = "REPLACE_WITH_YOUR_ZONE_ID"

# ── Database ───────────────────────────────────────────────────────────────────
db_instance_class = "db.t4g.medium"
db_name           = "digimap"
db_username       = "digimap"

# ── Cache ──────────────────────────────────────────────────────────────────────
redis_node_type = "cache.t4g.small"

# ── ECS ────────────────────────────────────────────────────────────────────────
# Update app_image before each deploy:
#   app_image = "123456789012.dkr.ecr.ap-southeast-1.amazonaws.com/digimap-backend:<git-sha>"
app_image         = "REPLACE_WITH_ECR_IMAGE_URI"
app_cpu           = "512"
app_memory        = "1024"
app_desired_count = 2

# ── App secrets (set via TF_VAR_ env var or -var flag — never commit) ──────────
# jwt_secret = "REPLACE_WITH_STRONG_SECRET"

# ── App config ─────────────────────────────────────────────────────────────────
jwt_issuer         = "digimap"
jwt_access_expiry  = "15m"
jwt_refresh_expiry = "7d"
