# Terraform — digimap-backend AWS Infrastructure

Manages the AWS infrastructure for `digimap-backend` across two environments: **develop** and **production**.

## Architecture

```
Route53 A alias
    └─► ALB (ACM TLS :443, HTTP→HTTPS redirect)
            ├─► Target Group blue  ─┐
            └─► Target Group green ─┴─► ECS Fargate tasks (:8080)
                                           ├─► RDS PostgreSQL 17 (private subnets)
                                           └─► ElastiCache Redis 7 (private subnets)

ECR ─ app image
S3  ─ main / assets / snapshot buckets
Secrets Manager ─ DB password, JWT secret, Redis auth token
IAM task role   ─ S3 access (replaces static AWS keys)
CodeDeploy      ─ blue/green deployments with auto-rollback
```

## Structure

```
terraform/
├── bootstrap/          # Remote state (run once)
├── modules/
│   ├── network/        # VPC, subnets, NAT
│   ├── security/       # Security groups
│   ├── ecr/            # ECR repository
│   ├── secrets/        # Secrets Manager
│   ├── database/       # RDS PostgreSQL 17
│   ├── cache/          # ElastiCache Redis 7
│   ├── storage/        # S3 buckets (main/assets/snapshot)
│   ├── alb/            # ALB, ACM, target groups
│   ├── iam/            # ECS execution + task roles
│   └── ecs-service/    # ECS cluster, service, CodeDeploy
└── envs/
    ├── develop/
    └── production/
```

## Prerequisites

- Terraform >= 1.10.0 (required for S3 native state locking)
- AWS CLI configured (`aws configure` or env vars)
- A Route53 hosted zone for your domain

## Step 1 — Bootstrap Remote State (once)

```bash
cd terraform/bootstrap
terraform init
terraform apply
# Note the output: state_bucket_name
```

The bucket name defaults to `digimap-terraform-state`. State locking uses S3's native lock file (`use_lockfile = true`) — no DynamoDB table required. If you change the bucket name, update `backend.tf` in each env.

## Step 2 — Set Required Variables

Set secrets via environment variables (never commit them):

```bash
export TF_VAR_jwt_secret="$(openssl rand -base64 48)"
export TF_VAR_route53_zone_id="YOUR_ZONE_ID"
```

Update the `*.tfvars` file for the target environment:
- Set `domain_name` (e.g. `api-dev.digimap.ai`)
- Set `route53_zone_id`
- Set `app_image` (set to your ECR URI with a placeholder tag for the initial apply; update per deploy)

## Step 3 — Apply an Environment

```bash
cd terraform/envs/develop       # or production
terraform init
terraform plan -var-file=develop.tfvars
terraform apply -var-file=develop.tfvars
```

Note the outputs — you'll need `ecr_repository_url`, `ecs_cluster_name`, `ecs_service_name`, `codedeploy_app_name`, and `codedeploy_deployment_group_name`.

## Step 4 — First Deploy

### 4a. Build and push the image

```bash
# Authenticate to ECR
aws ecr get-login-password --region ap-southeast-1 \
  | docker login --username AWS --password-stdin <ecr_repository_url>

# Build and push
GIT_SHA=$(git rev-parse --short HEAD)
docker build -t digimap-backend:${GIT_SHA} .
docker tag digimap-backend:${GIT_SHA} <ecr_repository_url>:${GIT_SHA}
docker push <ecr_repository_url>:${GIT_SHA}
```

### 4b. Run database migrations

```bash
aws ecs run-task \
  --cluster <ecs_cluster_name> \
  --task-definition digimap-<env>-app \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[<private-subnet-id>],securityGroups=[<app-sg-id>],assignPublicIp=DISABLED}" \
  --overrides '{"containerOverrides":[{"name":"app","command":["migrate","up"]}]}' \
  --region ap-southeast-1
```

Wait for the task to exit 0 before proceeding.

### 4c. Update `app_image` in the tfvars and re-apply

```bash
# In develop.tfvars:
app_image = "<ecr_repository_url>:<GIT_SHA>"

terraform apply -var-file=develop.tfvars
```

This updates the ECS task definition. The ECS service will pick it up on the next CodeDeploy deployment.

### 4d. Trigger a CodeDeploy deployment

Create an `appspec.yaml` for each deployment (substitute values from Terraform outputs):

```yaml
version: 0.0
Resources:
  - TargetService:
      Type: AWS::ECS::Service
      Properties:
        TaskDefinition: "<task_definition_arn>"
        LoadBalancerInfo:
          ContainerName: "app"
          ContainerPort: 8080
```

Then deploy:

```bash
aws deploy create-deployment \
  --application-name <codedeploy_app_name> \
  --deployment-group-name <codedeploy_deployment_group_name> \
  --revision '{"revisionType":"AppSpecContent","appSpecContent":{"content":"<appspec_yaml_as_json_string>"}}' \
  --region ap-southeast-1
```

Or use the AWS console: **CodeDeploy → Applications → digimap-\<env\>-app → Create deployment**.

CodeDeploy will:
1. Launch new tasks on the idle target group.
2. Wait for health checks to pass on `/health`.
3. Shift ALB traffic to the new target group.
4. Drain and terminate old tasks after 5 minutes.
5. Auto-rollback if health checks fail.

## Promote to Production

```bash
cd terraform/envs/production
terraform init
terraform plan -var-file=production.tfvars
terraform apply -var-file=production.tfvars
```

Production differences vs develop:
- `multi_az = true` on RDS
- 2-node Redis with automatic failover
- `desired_count = 2` ECS tasks
- `deletion_protection = true` on RDS
- Single NAT gateway → one NAT per AZ
- Larger instance sizes (`db.t4g.medium`, `cache.t4g.small`, 512 CPU / 1024 MB)
- Log retention 90 days (vs 30)

## Remove Static AWS Keys from the App

The app currently uses `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY` env vars
(`internal/platform/storage/s3.go`). With this Terraform setup the ECS task
role grants S3 access directly. To use it:

Change `s3.go` to use the default credential chain instead of a static provider:

```go
// Before:
cfg, err := config.LoadDefaultConfig(ctx,
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
        accessKey, secretKey, "",
    )),
)

// After:
cfg, err := config.LoadDefaultConfig(ctx) // uses task role automatically
```

The `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY` env vars can then be removed
from the task definition (they're not set by this Terraform config).

## Extensions (not yet provisioned)

These are coded in the app but behind Log stubs. Add their Terraform modules
when the Go code is wired up:

| Service | Module to add | Config needed in `config.Load()` |
|---|---|---|
| OpenSearch | `modules/opensearch/` | `OPENSEARCH_ENDPOINT/USER/PASSWORD` |
| SES email | SES domain identity + DKIM | `SES_FROM_EMAIL` domain verification |
| CloudFront | `modules/cloudfront/` | `AWS_CF_DISTRIBUTION_ID` |
| Firebase FCM | External (Google) | `FIREBASE_CREDENTIALS_PATH` |
