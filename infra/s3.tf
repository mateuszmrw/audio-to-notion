resource "aws_s3_bucket" "audio-to-notion-bucket" {
  bucket = "audio-to-notion-bucket"
}

resource "aws_s3_bucket_cors_configuration" "s3-cors-audio-to-notion-bucket" {
  bucket = aws_s3_bucket.audio-to-notion-bucket.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket_public_access_block" "audio-to-notion-public-access-block" {
  bucket = aws_s3_bucket.audio-to-notion-bucket.id

  block_public_acls       = true
  block_public_policy     = true
  restrict_public_buckets = true
  ignore_public_acls      = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "audio-to-notion-bucket-encryption" {
  bucket = aws_s3_bucket.audio-to-notion-bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_lambda_permission" "allow-audio-to-notion-invoke" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.audio-to-notion-lambda.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.audio-to-notion-bucket.arn
}

resource "aws_s3_bucket_notification" "audio-to-notion-bucket-notification" {
  bucket = aws_s3_bucket.audio-to-notion-bucket.bucket

  lambda_function {
    lambda_function_arn = aws_lambda_function.audio-to-notion-lambda.arn
    events              = ["s3:ObjectCreated:*"]
  }
}