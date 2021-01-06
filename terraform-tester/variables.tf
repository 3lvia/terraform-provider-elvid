variable "tenant_id" {
}

variable "terraform_sp_client_id" {
}

variable "terraform_sp_client_secret" {
}

variable "environment" {
  default = "dev"
}

variable "system_name" {
  default = "examplesystem"
}

variable "domains" {
  type = map
  default = {
    "dev" = ["http://localhost:{port}", "https://examplesystem.dev-elvia.io"]
    "test" = ["https://examplesystem.test-elvia.io"]
    "prod" = ["https://examplesystem.elvia.io"]
  }
}