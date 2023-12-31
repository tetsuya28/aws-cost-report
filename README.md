## What is aws_cost_usage
Notify daily cost usage of an AWS account to slack channel.

## How to use this
You can deploy with Terraform resources on your AWS account.

```hcl
module "cost" {
  source        = "github.com/tetsuya28/aws_cost_report.git//module"
  name          = "aws-cost-report" # AWS resource name to deploy Lambda, IAM, etc
  slack_channel = "#random"         # Slack channel to notify
  slack_token   = "xoxb-xxx"        # Slack token
  build_version = "v0.1.0"          # default : latest
}
```

## Development
- Init
```
cp .env{.sample,}
```

- Run locally
```
make run
```
