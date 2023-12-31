terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    github = {
      source = "integrations/github"
    }
  }
}

provider "github" {
  owner = "tetsuya28"
}
