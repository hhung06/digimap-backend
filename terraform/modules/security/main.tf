terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  required_version = "~> 1.9"
}

locals {
  common_tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ── ALB Security Group ────────────────────────────────────────────────────────

resource "aws_security_group" "alb" {
  name        = "${var.project}-${var.env}-alb"
  description = "Allow HTTP/HTTPS inbound to the Application Load Balancer"
  vpc_id      = var.vpc_id

  tags = merge(local.common_tags, {
    Name = "${var.project}-${var.env}-alb"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "alb_ingress_http_ipv4" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "HTTP from anywhere (IPv4)"
}

resource "aws_security_group_rule" "alb_ingress_http_ipv6" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80
  ipv6_cidr_blocks  = ["::/0"]
  description       = "HTTP from anywhere (IPv6)"
}

resource "aws_security_group_rule" "alb_ingress_https_ipv4" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "HTTPS from anywhere (IPv4)"
}

resource "aws_security_group_rule" "alb_ingress_https_ipv6" {
  security_group_id = aws_security_group.alb.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 443
  to_port           = 443
  ipv6_cidr_blocks  = ["::/0"]
  description       = "HTTPS from anywhere (IPv6)"
}

resource "aws_security_group_rule" "alb_egress_all" {
  security_group_id = aws_security_group.alb.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "Allow all outbound"
}

# ── App Security Group ────────────────────────────────────────────────────────

resource "aws_security_group" "app" {
  name        = "${var.project}-${var.env}-app"
  description = "Allow inbound from ALB on app port; allow all outbound"
  vpc_id      = var.vpc_id

  tags = merge(local.common_tags, {
    Name = "${var.project}-${var.env}-app"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "app_ingress_from_alb" {
  security_group_id        = aws_security_group.app.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 8080
  to_port                  = 8080
  source_security_group_id = aws_security_group.alb.id
  description              = "App port from ALB"
}

resource "aws_security_group_rule" "app_egress_all" {
  security_group_id = aws_security_group.app.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "Allow all outbound (RDS, Redis, AWS APIs)"
}

# ── RDS Security Group ────────────────────────────────────────────────────────

resource "aws_security_group" "rds" {
  name        = "${var.project}-${var.env}-rds"
  description = "Allow PostgreSQL inbound from app tier only"
  vpc_id      = var.vpc_id

  tags = merge(local.common_tags, {
    Name = "${var.project}-${var.env}-rds"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "rds_ingress_from_app" {
  security_group_id        = aws_security_group.rds.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 5432
  to_port                  = 5432
  source_security_group_id = aws_security_group.app.id
  description              = "PostgreSQL from app tier"
}

resource "aws_security_group_rule" "rds_egress_all" {
  security_group_id = aws_security_group.rds.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "Allow all outbound"
}

# ── Redis Security Group ──────────────────────────────────────────────────────

resource "aws_security_group" "redis" {
  name        = "${var.project}-${var.env}-redis"
  description = "Allow Redis inbound from app tier only"
  vpc_id      = var.vpc_id

  tags = merge(local.common_tags, {
    Name = "${var.project}-${var.env}-redis"
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_security_group_rule" "redis_ingress_from_app" {
  security_group_id        = aws_security_group.redis.id
  type                     = "ingress"
  protocol                 = "tcp"
  from_port                = 6379
  to_port                  = 6379
  source_security_group_id = aws_security_group.app.id
  description              = "Redis from app tier"
}

resource "aws_security_group_rule" "redis_egress_all" {
  security_group_id = aws_security_group.redis.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  description       = "Allow all outbound"
}
