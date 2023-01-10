# Used by terraform-tester.tf
terraform {
  required_version = ">= 1.0"
  required_providers {
    elvid = {
       source = "3lvia/elvid"
    }
  }
}
