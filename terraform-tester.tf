# This is used to test the provider. Check the readme for info about installing and running this

# Provider
provider "elvid" {
  tenant_id                  = var.tenant_id
  terraform_sp_client_id     = var.terraform_sp_client_id
  terraform_sp_client_secret = var.terraform_sp_client_secret
  environment                = var.environment
  override_elvid_authority   = "https://localhost:44383"
  #override_elvid_authority = "https://elvid.dev-elvia.io"
  run_hashed_secret_validation = true
}

provider "vault" {
  address = "https://vault.dev-elvia.io"
  auth_login {
    path = "auth/approle/login"

    parameters = {
      role_id = var.elvid_dev_vault_role_id
    }
  }
}

## User client

# resource "elvid_userclient" "userclient" {
#   client_name                            = "test"
#   scopes                                 = ["louvre.imageapi.useraccess", "profile", "openid", "ad_groups"]
#   domains                                = var.domains[var.environment]
#   redirect_uri_paths                     = ["/callback.html"]
#   post_logout_redirect_uri_paths         = ["/index.htm"]
#   local_login_enabled                    = true
#   idporten_login_enabled                 = true
#   elvia_ad_login_enabled                 = true
#   test_user_login_enabled                = false
#   require_client_secret                  = false
#   access_token_life_time                 = 3591
#   always_include_user_claims_in_id_token = true
#   client_name_language_key               = null
#   allow_use_of_refresh_tokens            = false
#   one_time_usage_for_refresh_tokens      = true
#   refresh_token_life_time                = 2592000
#   resource_taint_version = "2"
#   client_properties {
#     key   = "ad_groups_filter"
#     values = []
#   }
# }

# output "userclient" {
#   value = elvid_userclient.userclient
# }

## Machine client

# resource "elvid_machineclient" "machineclient10" {
#   name                    = "2024-09-23"
#   test_user_login_enabled = true
#   access_token_life_time  = 3522
#   scopes                  = ["elvid.verifydeployment", "louvre.imageapi"]
#   resource_taint_version  = "5"
#   # client_claims {
#   #   type   = "client_dna_topics_read12"
#   #   values = ["topicA", "topicB", "C"]
#   # }
#   # client_claims {
#   #   type   = "client_edna_topics_write"
#   #   values = ["topicA", "topicB", "D"]
#   # }
# }

# resource "elvid_clientsecret" "clientsecret" {
#   client_id              = elvid_machineclient.machineclient10.id
# }

# output "machineclient" {
#   value = elvid_machineclient.machineclient10.client_id
# }

# output "clientsecret" {
#   value = nonsensitive(elvid_clientsecret.clientsecret.secret_value)
# }

## API scope
# resource "elvid_apiscope" "apiscope" {
#     name = "terraform-provider-elvid-tester-apiscope"
#     description = "Scope opprettet fra test av Elvid Terraform provider (terraform-tester i terraform-provider-elvid)"
#     user_claims = ["email", "ad_groups"]
#     allow_user_clients = true
#     resource_taint_version = "3"
# }

## Module userclient
# module "elvid_userclient" {
#   source  = "app.terraform.io/Elvia/userclient/elvid"
#   # source      = "C:\\3lvia\\terraform-elvid-userclient"
#   environment = "dev"
#   client_name = "test-userclient"
#   scopes = [ "louvre.imageapi.useraccess", "openid", "ad_groups"]
#   domains = var.domains[var.environment]
#   redirect_uri_paths = [ "/silentcallback.html", "/oidc/callback"]
#   post_logout_redirect_uri_paths = [""]
#   elvia_ad_login_enabled         = true
#   system_name      = "elvid"
#   client_secret_enabled = false
#   ad_groups_filter = ["test"]
#   allow_use_of_refresh_tokens = false
#   one_time_usage_for_refresh_tokens  = false
# }

## Module machineclient
# module "elvid_machineclient" {
#   # source  = "app.terraform.io/Elvia/machineclient/elvid"
#   source      = "C:\\3lvia\\terraform-elvid-machineclient"
#   scopes = ["louvre.imageapi"]
#   environment      = var.environment
#   system_name      = "elvid"
#   application_name = "test-machineclient"
#   client_claims = [
#     {
#       type = "client_kafka_topic_read"
#       values = ["topic1", "topic2"]
#     },
#     {
#       type = "client_kafka_topic_write"
#       values = ["topic1", "topic2"]
#     }
#   ]
# }

variable "tenant_id" {
}

variable "terraform_sp_client_id" {
}

variable "terraform_sp_client_secret" {
}

variable "elvid_dev_vault_role_id" {
}


variable "environment" {
  default = "dev"
}

variable "system_name" {
  default = "examplesystem"
}

variable "domains" {
  type = map(any)
  default = {
    "dev"  = ["http://localhost:{port}", "https://examplesystem.dev-elvia.io"]
    "test" = ["https://examplesystem.test-elvia.io"]
    "prod" = ["https://examplesystem.elvia.io"]
  }
}
