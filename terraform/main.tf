locals {
  name       = "aws-cost-usage"
  created_by = "Terraform"
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      CreatedBy = local.created_by
    }
  }
}

terraform {
  backend "s3" {
    key    = "aws-cost-usage/default.tfstate"
    region = "ap-northeast-1"
  }
}

resource "aws_ecrpublic_repository" "foo" {
  repository_name = local.name
}
