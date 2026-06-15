terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

locals {
  name_prefix = "${var.project}-${var.env}"
}

# ── Latest Amazon Linux 2023 AMI ──────────────────────────────────────────────

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# ── Security Group ────────────────────────────────────────────────────────────

resource "aws_security_group" "ec2" {
  name        = "${local.name_prefix}-ec2"
  description = "EC2 app instance — dev"
  vpc_id      = var.vpc_id

  dynamic "ingress" {
    for_each = length(var.allowed_ssh_cidrs) > 0 ? [1] : []
    content {
      description = "SSH"
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      cidr_blocks = var.allowed_ssh_cidrs
    }
  }

  ingress {
    description = "App port"
    from_port   = var.app_port
    to_port     = var.app_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name    = "${local.name_prefix}-ec2"
    Project = var.project
    Env     = var.env
  }
}

# ── IAM Instance Profile ──────────────────────────────────────────────────────

data "aws_iam_policy_document" "ec2_assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ec2" {
  name               = "${local.name_prefix}-ec2"
  assume_role_policy = data.aws_iam_policy_document.ec2_assume.json

  tags = {
    Project = var.project
    Env     = var.env
  }
}

data "aws_iam_policy_document" "ec2_permissions" {
  statement {
    sid     = "ReadSecrets"
    actions = ["secretsmanager:GetSecretValue"]
    resources = [
      var.db_password_secret_arn,
      var.jwt_secret_arn,
      var.redis_password_secret_arn,
    ]
  }

  statement {
    sid = "S3Objects"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [for arn in var.s3_bucket_arns : "${arn}/*"]
  }

  statement {
    sid       = "S3List"
    actions   = ["s3:ListBucket"]
    resources = var.s3_bucket_arns
  }

  statement {
    sid = "ECRAuth"
    actions = [
      "ecr:GetAuthorizationToken",
    ]
    resources = ["*"]
  }

  statement {
    sid = "ECRPull"
    actions = [
      "ecr:BatchGetImage",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchCheckLayerAvailability",
    ]
    resources = ["arn:aws:ecr:${var.aws_region}:*:repository/*"]
  }

  statement {
    sid = "CloudWatchLogs"
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "ec2" {
  name   = "${local.name_prefix}-ec2"
  role   = aws_iam_role.ec2.id
  policy = data.aws_iam_policy_document.ec2_permissions.json
}

resource "aws_iam_instance_profile" "ec2" {
  name = "${local.name_prefix}-ec2"
  role = aws_iam_role.ec2.name
}

# ── EC2 Instance ──────────────────────────────────────────────────────────────

resource "aws_instance" "this" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = var.instance_type
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2.name
  key_name               = var.key_name != "" ? var.key_name : null

  user_data = base64encode(templatefile("${path.module}/userdata.sh.tpl", {
    region                    = var.aws_region
    ecr_registry              = var.ecr_registry
    app_image                 = var.app_image
    app_port                  = var.app_port
    env                       = var.env
    db_host                   = var.db_host
    db_port                   = var.db_port
    db_name                   = var.db_name
    db_username               = var.db_username
    db_sslmode                = var.db_sslmode
    redis_addr                = var.redis_addr
    s3_bucket                 = var.s3_bucket
    s3_assets_bucket          = var.s3_assets_bucket
    s3_snapshot_bucket        = var.s3_snapshot_bucket
    jwt_issuer                = var.jwt_issuer
    jwt_access_expiry         = var.jwt_access_expiry
    jwt_refresh_expiry        = var.jwt_refresh_expiry
    db_password_secret_arn    = var.db_password_secret_arn
    jwt_secret_arn            = var.jwt_secret_arn
    redis_password_secret_arn = var.redis_password_secret_arn
  }))

  root_block_device {
    volume_type           = "gp3"
    volume_size           = 20
    delete_on_termination = true
  }

  tags = {
    Name    = "${local.name_prefix}-app"
    Project = var.project
    Env     = var.env
  }
}

# ── Elastic IP ────────────────────────────────────────────────────────────────

resource "aws_eip" "this" {
  instance = aws_instance.this.id
  domain   = "vpc"

  tags = {
    Name    = "${local.name_prefix}-app"
    Project = var.project
    Env     = var.env
  }
}
