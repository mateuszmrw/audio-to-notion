resource "aws_cloudwatch_log_group" "audio-to-notion-lambda-log-group" {
  name              = "/aws/lambda/${aws_lambda_function.audio-to-notion-lambda.function_name}"
  retention_in_days = 7
}