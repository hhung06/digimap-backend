terraform {
  backend "s3" {
    bucket       = "digimap-terraform-state"
    key          = "production/terraform.tfstate"
    region       = "ap-northeast-1"
    encrypt      = true
    use_lockfile = true # S3 native locking (Terraform >= 1.10, no DynamoDB needed)
  }
}
