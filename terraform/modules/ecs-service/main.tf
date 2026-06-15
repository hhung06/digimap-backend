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

# ─── CloudWatch Log Group ─────────────────────────────────────────────────────

resource "aws_cloudwatch_log_group" "this" {
  name              = "/ecs/${var.project}/${var.env}"
  retention_in_days = var.log_retention_days

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ─── ECS Cluster ─────────────────────────────────────────────────────────────

resource "aws_ecs_cluster" "this" {
  name = local.name_prefix

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ─── ECS Task Definition ─────────────────────────────────────────────────────

resource "aws_ecs_task_definition" "this" {
  family                   = "${local.name_prefix}-app"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.cpu
  memory                   = var.memory
  execution_role_arn       = var.execution_role_arn
  task_role_arn            = var.task_role_arn

  container_definitions = jsonencode([
    {
      name      = "app"
      image     = var.app_image
      essential = true

      portMappings = [
        {
          containerPort = var.app_port
          protocol      = "tcp"
        }
      ]

      environment = [
        { name = "PORT", value = tostring(var.app_port) },
        { name = "ENV", value = var.env },
        { name = "LOG_LEVEL", value = var.log_level },
        { name = "LOG_FORMAT", value = var.log_format },
        { name = "DB_HOST", value = var.db_host },
        { name = "DB_PORT", value = var.db_port },
        { name = "DB_NAME", value = var.db_name },
        { name = "DB_USER", value = var.db_username },
        { name = "DB_SSLMODE", value = var.db_sslmode },
        { name = "DB_MAX_CONNS", value = var.db_max_conns },
        { name = "DB_MIN_CONNS", value = var.db_min_conns },
        { name = "REDIS_ADDR", value = var.redis_addr },
        { name = "AWS_S3_BUCKET", value = var.s3_bucket },
        { name = "AWS_S3_ASSETS_BUCKET", value = var.s3_assets_bucket },
        { name = "AWS_S3_SNAPSHOT_BUCKET", value = var.s3_snapshot_bucket },
        { name = "AWS_REGION", value = var.aws_region },
        { name = "JWT_ISSUER", value = var.jwt_issuer },
        { name = "JWT_ACCESS_EXPIRY", value = var.jwt_access_expiry },
        { name = "JWT_REFRESH_EXPIRY", value = var.jwt_refresh_expiry },
      ]

      secrets = [
        {
          name      = "DB_PASSWORD"
          valueFrom = var.db_password_secret_arn
        },
        {
          name      = "JWT_SECRET_KEY"
          valueFrom = var.jwt_secret_arn
        },
        {
          name      = "REDIS_PASSWORD"
          valueFrom = var.redis_password_secret_arn
        },
      ]

      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = "/ecs/${var.project}/${var.env}"
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "app"
        }
      }

      healthCheck = {
        command     = ["CMD-SHELL", "wget --quiet --tries=1 --spider http://localhost:${var.app_port}/health || exit 1"]
        interval    = 30
        timeout     = 10
        retries     = 3
        startPeriod = 60
      }
    }
  ])

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ─── ECS Service ─────────────────────────────────────────────────────────────

resource "aws_ecs_service" "this" {
  name            = "${local.name_prefix}-app"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.this.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.target_group_blue_arn
    container_name   = "app"
    container_port   = var.app_port
  }

  deployment_controller {
    type = "CODE_DEPLOY"
  }

  health_check_grace_period_seconds = 120

  # CodeDeploy manages task_definition and load_balancer after initial deployment
  lifecycle {
    ignore_changes = [task_definition, load_balancer]
  }

  depends_on = [aws_cloudwatch_log_group.this]

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ─── CodeDeploy Application ───────────────────────────────────────────────────

resource "aws_codedeploy_app" "this" {
  name             = "${local.name_prefix}-app"
  compute_platform = "ECS"

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

# ─── CodeDeploy IAM Role ──────────────────────────────────────────────────────

data "aws_iam_policy_document" "codedeploy_assume_role" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["codedeploy.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "codedeploy" {
  name               = "${local.name_prefix}-codedeploy"
  assume_role_policy = data.aws_iam_policy_document.codedeploy_assume_role.json

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "codedeploy_ecs" {
  role       = aws_iam_role.codedeploy.name
  policy_arn = "arn:aws:iam::aws:policy/AWSCodeDeployRoleForECS"
}

# ─── CodeDeploy Deployment Group ──────────────────────────────────────────────

resource "aws_codedeploy_deployment_group" "this" {
  app_name               = aws_codedeploy_app.this.name
  deployment_group_name  = "${local.name_prefix}-app"
  service_role_arn       = aws_iam_role.codedeploy.arn
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"]
  }

  blue_green_deployment_config {
    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }

    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }
  }

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  ecs_service {
    cluster_name = aws_ecs_cluster.this.name
    service_name = aws_ecs_service.this.name
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [var.https_listener_arn]
      }

      target_group {
        name = var.target_group_blue_name
      }

      target_group {
        name = var.target_group_green_name
      }
    }
  }

  tags = {
    Project   = var.project
    Env       = var.env
    ManagedBy = "terraform"
  }
}
