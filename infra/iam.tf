data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "audio-to-notion-lambda-role" {
  name               = "AssumeLambdaRole"
  description        = "Role for lambda to assume lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

data "aws_iam_policy_document" "audio-to-notion-lambda-logging-policy-document" {
  statement {
    effect = "Allow"
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "arn:aws:logs:*:*:*",
    ]
  }
}

resource "aws_iam_policy" "audio-to-notion-lambda-logging-policy" {
  name        = "AllowLambdaLoggingPolicy"
  description = "Policy for lambda cloudwatch logging"
  policy      = data.aws_iam_policy_document.audio-to-notion-lambda-logging-policy-document.json
}

resource "aws_iam_role_policy_attachment" "audio-to-notion-lambda-logging-policy-attachment" {
  role       = aws_iam_role.audio-to-notion-lambda-role.id
  policy_arn = aws_iam_policy.audio-to-notion-lambda-logging-policy.arn
}

data "aws_iam_policy_document" "audio-to-notion-s3-get-object-policy-document" {
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject"
    ]

    resources = [aws_s3_bucket.audio-to-notion-bucket.arn, "${aws_s3_bucket.audio-to-notion-bucket.arn}/*"]
  }
}

resource "aws_iam_policy" "audio-to-notion-lambda-s3-get-object-policy" {
  name   = "AllowAudioNotionGetObjectPolicy"
  policy = data.aws_iam_policy_document.audio-to-notion-s3-get-object-policy-document.json
}

resource "aws_iam_role_policy_attachment" "audio-to-notion-lambda-s3-get-object-policy-attachment" {
  role       = aws_iam_role.audio-to-notion-lambda-role.id
  policy_arn = aws_iam_policy.audio-to-notion-lambda-s3-get-object-policy.arn
}

