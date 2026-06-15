terraform {
  backend "s3" {
    # Fill in after running terraform/bootstrap:
    #   bucket = "<output: state_bucket_name>"
    bucket       = "digimap-terraform-state"
    key          = "develop/terraform.tfstate"
    region       = "ap-northeast-1"
    encrypt      = true
    use_lockfile = true # S3 native locking (Terraform >= 1.10, no DynamoDB needed)
  }
}
