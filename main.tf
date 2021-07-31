variable "profile" {}
variable "region" {}
variable "slack_token" {}
variable "slack_channel" {}
variable "schedule_expression" {
  default = "cron(0 0 * * ? *)"
}

locals {
  name       = "aws_cost_usage"
  created_by = "Terraform"
}

provider "aws" {
  region  = var.region
  profile = var.profile
  default_tags {
    tags = {
      App       = local.name
      CreatedBy = local.created_by
    }
  }
}

terraform {
  backend "s3" {
    key    = "aws_cost_usage.tfstate"
    region = "ap-northeast-1"
  }
}

data "archive_file" "this" {
  type        = "zip"
  source_file = "bin/main"
  output_path = "bin/main.zip"
}

resource "aws_lambda_function" "this" {
  function_name    = local.name
  filename         = "bin/main.zip"
  handler          = "main"
  source_code_hash = data.archive_file.this.output_base64sha256
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 10
  role             = aws_iam_role.this.arn
  environment {
    variables = {
      "SLACK_TOKEN"   = var.slack_token
      "SLACK_CHANNEL" = var.slack_channel
    }
  }
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${local.name}"
  retention_in_days = 7
}

resource "aws_cloudwatch_event_rule" "this" {
  name                = local.name
  schedule_expression = var.schedule_expression
}

resource "aws_cloudwatch_event_target" "this" {
  rule      = aws_cloudwatch_event_rule.this.name
  target_id = local.name
  arn       = aws_lambda_function.this.arn
}

resource "aws_iam_role" "this" {
  name               = local.name
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
  name        = local.name
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
