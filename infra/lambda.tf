
locals {
  binary_path  = "${path.module}/lambda/audio-to-notion/bin/main"
  src_path     = "${path.module}/lambda/audio-to-notion/main.go"
  archive_path = "${path.module}/lambda/main.zip"
  binary_name  = "main"
}


resource "null_resource" "function-binary" {
  provisioner "local-exec" {
    command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
  }
  triggers = {
    always_run = timestamp()
  }
}

data "archive_file" "lambda-archive" {
  depends_on = [null_resource.function-binary]

  type        = "zip"
  source_file = local.binary_path
  output_path = local.archive_path
}

resource "aws_lambda_function" "audio-to-notion-lambda" {
  function_name    = "audio-to-notion"
  role             = aws_iam_role.audio-to-notion-lambda-role.arn
  handler          = local.binary_name
  memory_size      = 128
  timeout          = 300
  filename         = local.archive_path
  source_code_hash = data.archive_file.lambda-archive.output_base64sha256

  environment {
    variables = {
      OPEN_AI_KEY        = var.OPEN_AI_KEY
      NOTION_TOKEN       = var.NOTION_TOKEN
      NOTION_DATABASE_ID = var.NOTION_DATABASE_ID
    }
  }

  runtime = "provided.al2"
}