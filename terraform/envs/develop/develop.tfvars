# ── Core ───────────────────────────────────────────────────────────────────────
aws_region = "ap-northeast-1"
project    = "digimap"
env        = "develop"

# ── Network ────────────────────────────────────────────────────────────────────
vpc_cidr             = "10.0.0.0/16"
azs                  = ["ap-northeast-1a", "ap-northeast-1c"]
public_subnet_cidrs  = ["10.0.1.0/24", "10.0.2.0/24"]
private_subnet_cidrs = ["10.0.11.0/24", "10.0.12.0/24"]

# ── Database ───────────────────────────────────────────────────────────────────
db_instance_class = "db.t4g.micro"
db_name           = "digimap"
db_username       = "digimap"

# ── Cache ──────────────────────────────────────────────────────────────────────
redis_node_type = "cache.t4g.micro"

# ── EC2 ────────────────────────────────────────────────────────────────────────
instance_type = "t3.small"
key_name      = ""   # set to your EC2 key pair name to enable SSH

# Set to your public IP to restrict SSH access, e.g. ["1.2.3.4/32"]
allowed_ssh_cidrs = []

# Update before each deploy:
#   app_image = "123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/digimap-backend:<git-sha>"
app_image = "REPLACE_WITH_ECR_IMAGE_URI"

# ── App secrets (set via TF_VAR_ env var — never commit real values) ───────────
# jwt_secret = "REPLACE_WITH_STRONG_SECRET"

# ── App config ─────────────────────────────────────────────────────────────────
jwt_issuer         = "digimap-develop"
jwt_access_expiry  = "15m"
jwt_refresh_expiry = "7d"
