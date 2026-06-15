output "main_bucket_name" {
  description = "Name of the main uploads bucket (AWS_S3_BUCKET)"
  value       = aws_s3_bucket.this["main"].id
}

output "assets_bucket_name" {
  description = "Name of the assets bucket (AWS_S3_ASSETS_BUCKET)"
  value       = aws_s3_bucket.this["assets"].id
}

output "snapshot_bucket_name" {
  description = "Name of the snapshot bucket (AWS_S3_SNAPSHOT_BUCKET)"
  value       = aws_s3_bucket.this["snapshot"].id
}

output "main_bucket_arn" {
  description = "ARN of the main uploads bucket"
  value       = aws_s3_bucket.this["main"].arn
}

output "assets_bucket_arn" {
  description = "ARN of the assets bucket"
  value       = aws_s3_bucket.this["assets"].arn
}

output "snapshot_bucket_arn" {
  description = "ARN of the snapshot bucket"
  value       = aws_s3_bucket.this["snapshot"].arn
}

output "bucket_arns" {
  description = "List of all three bucket ARNs (for use in IAM policies)"
  value = [
    aws_s3_bucket.this["main"].arn,
    aws_s3_bucket.this["assets"].arn,
    aws_s3_bucket.this["snapshot"].arn,
  ]
}
