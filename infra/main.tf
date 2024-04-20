provider "aws" {
  region  = "eu-central-1"
  profile = "flashcard"
}

terraform {
  backend "s3" {
    bucket = "audio-to-notion-terraform"
    key    = "state"
    region = var.AWS_REGION
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}