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
  required_version = "~> 1.9"
}

resource "random_password" "db_password" {
  length           = 32
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# ---------------------------------------------------------------------------
# DB password secret
# ---------------------------------------------------------------------------

resource "aws_secretsmanager_secret" "db_password" {
  name                    = "/${var.project}/${var.env}/db-password"
  recovery_window_in_days = 7

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_secretsmanager_secret_version" "db_password" {
  secret_id     = aws_secretsmanager_secret.db_password.id
  secret_string = random_password.db_password.result
}

# ---------------------------------------------------------------------------
# JWT secret
# ---------------------------------------------------------------------------

resource "aws_secretsmanager_secret" "jwt_secret" {
  name                    = "/${var.project}/${var.env}/jwt-secret"
  recovery_window_in_days = 7

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_secretsmanager_secret_version" "jwt_secret" {
  secret_id     = aws_secretsmanager_secret.jwt_secret.id
  secret_string = var.jwt_secret
}

# ---------------------------------------------------------------------------
# IAM policy granting read access to both secrets
# ---------------------------------------------------------------------------

data "aws_iam_policy_document" "read_secrets" {
  statement {
    sid    = "AllowGetSecretValue"
    effect = "Allow"

    actions = ["secretsmanager:GetSecretValue"]

    resources = [
      aws_secretsmanager_secret.db_password.arn,
      aws_secretsmanager_secret.jwt_secret.arn,
    ]
  }
}

resource "aws_iam_policy" "read_secrets" {
  name        = "${var.project}-${var.env}-read-secrets"
  description = "Grants GetSecretValue on ${var.project}/${var.env} secrets"
  policy      = data.aws_iam_policy_document.read_secrets.json

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}
