variable "name" {}
variable "slack_token" {}
variable "slack_channel" {}
variable "schedule_expression" {
  default = "cron(0 0 * * ? *)"
}
variable "build_version" {
  default = "v0.0.0"
}

resource "aws_lambda_function" "this" {
  function_name = var.name
  s3_bucket     = "tetsuya28-aws-cost-report"
  s3_key        = "${var.build_version}/main.zip"
  runtime       = "go1.x"
  handler       = "main"
  memory_size   = 128
  timeout       = 10
  role          = aws_iam_role.this.arn
  environment {
    variables = {
      "SLACK_TOKEN"   = var.slack_token
      "SLACK_CHANNEL" = var.slack_channel
    }
  }
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${var.name}"
  retention_in_days = 7
}

resource "aws_cloudwatch_event_rule" "this" {
  name                = var.name
  schedule_expression = var.schedule_expression
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = aws_cloudwatch_event_rule.this.name
  target_id = var.name
  arn       = aws_lambda_function.this.arn
}

resource "aws_iam_role" "this" {
  name               = var.name
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "this" {
  name        = var.name
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy      = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents",
        "ce:*"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "this" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.this.arn
}

resource "aws_lambda_permission" "this" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.this.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.this.arn
}
