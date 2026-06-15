#!/bin/bash
set -euo pipefail

# ── Install Docker & AWS CLI ──────────────────────────────────────────────────
dnf update -y
dnf install -y docker
systemctl enable --now docker
usermod -aG docker ec2-user

# ── Create app directory ───────────────────────────────────────────────────────
mkdir -p /opt/digimap
chmod 700 /opt/digimap

# ── Fetch secrets from Secrets Manager ───────────────────────────────────────
DB_PASSWORD=$(aws secretsmanager get-secret-value \
  --secret-id "${db_password_secret_arn}" \
  --query SecretString --output text --region "${region}")

JWT_SECRET=$(aws secretsmanager get-secret-value \
  --secret-id "${jwt_secret_arn}" \
  --query SecretString --output text --region "${region}")

REDIS_PASSWORD=$(aws secretsmanager get-secret-value \
  --secret-id "${redis_password_secret_arn}" \
  --query SecretString --output text --region "${region}")

# ── Write env file ─────────────────────────────────────────────────────────────
cat > /opt/digimap/.env <<EOF
PORT=${app_port}
ENV=${env}
LOG_LEVEL=debug
LOG_FORMAT=text
DB_HOST=${db_host}
DB_PORT=${db_port}
DB_NAME=${db_name}
DB_USER=${db_username}
DB_PASSWORD=$DB_PASSWORD
DB_SSLMODE=${db_sslmode}
DB_MAX_CONNS=10
DB_MIN_CONNS=2
REDIS_ADDR=${redis_addr}
REDIS_PASSWORD=$REDIS_PASSWORD
AWS_S3_BUCKET=${s3_bucket}
AWS_S3_ASSETS_BUCKET=${s3_assets_bucket}
AWS_S3_SNAPSHOT_BUCKET=${s3_snapshot_bucket}
AWS_REGION=${region}
JWT_SECRET_KEY=$JWT_SECRET
JWT_ISSUER=${jwt_issuer}
JWT_ACCESS_EXPIRY=${jwt_access_expiry}
JWT_REFRESH_EXPIRY=${jwt_refresh_expiry}
EOF
chmod 600 /opt/digimap/.env

# ── Login to ECR ───────────────────────────────────────────────────────────────
aws ecr get-login-password --region "${region}" \
  | docker login --username AWS --password-stdin "${ecr_registry}"

# ── Pull image ─────────────────────────────────────────────────────────────────
docker pull "${app_image}"

# ── Systemd service ────────────────────────────────────────────────────────────
cat > /etc/systemd/system/digimap-app.service <<'UNIT'
[Unit]
Description=digimap-app
After=docker.service network-online.target
Requires=docker.service

[Service]
Restart=always
RestartSec=5
ExecStartPre=-/usr/bin/docker stop digimap-app
ExecStartPre=-/usr/bin/docker rm digimap-app
ExecStart=/usr/bin/docker run --name digimap-app \
  --env-file /opt/digimap/.env \
  -p ${app_port}:${app_port} \
  ${app_image} serve
ExecStop=/usr/bin/docker stop digimap-app

[Install]
WantedBy=multi-user.target
UNIT

systemctl daemon-reload
systemctl enable --now digimap-app
