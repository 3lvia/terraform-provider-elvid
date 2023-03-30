# This is used to test the provider. Check the readme for info about installing and running this

# Provider
provider "elvid" {
  tenant_id = var.tenant_id
  terraform_sp_client_id = var.terraform_sp_client_id
  terraform_sp_client_secret = var.terraform_sp_client_secret
  environment = var.environment
  override_elvid_authority = "https://localhost:44383"
  # override_elvid_authority = "https://elvid.dev-elvia.io"
}

provider "vault" {
  auth_login {
    path = "auth/approle/login"

    parameters = {
      role_id = var.elvid_dev_vault_role_id
    }
  }
}

## User client

resource "elvid_userclient" "userclient" {
    client_name = "test"
    scopes = ["openid", "ad_groups"]
    domains = var.domains[var.environment]
    redirect_uri_paths = ["/callback.html"]
    post_logout_redirect_uri_paths = ["/index.htm"]
    bankid_login_enabled = true
    local_login_enabled = true
    elvia_ad_login_enabled = true
    test_user_login_enabled = false
    require_client_secret = false
    access_token_life_time = 3598
    always_include_user_claims_in_id_token = true
    client_name_language_key = null
    allow_use_of_refresh_tokens = false
    one_time_usage_for_refresh_tokens = true
    refresh_token_life_time = 2592000
}

# output "userclient" {
#   value = elvid_userclient.userclient
# }

## Machine client

# resource "elvid_machineclient" "machineclient" {
#     name = "onsdag"
#     test_user_login_enabled = true
#     access_token_life_time = 3511
#     resource_taint_version = "1"
#     scopes = ["elvid.verifydeployment"]
#     client_claims {
#       type = "edna_topics_read"
#       values = ["topic1"]
#     }
#     client_claims {
#       type = "edna_topics_write"
#       values = ["topicA", "topicB", "D"]
#     }
# }

# resource "elvid_clientsecret" "clientsecret" {
#     client_id = elvid_machineclient.machineclient.id
#     resource_taint_version = "2"
# }

# output "machineclient" {
#   value = elvid_machineclient.machineclient
# }

# output "clientsecret" {
#   value = elvid_clientsecret.clientsecret
# }

## API scope
# resource "elvid_apiscope" "apiscope" {
#     name = "terraform-provider-elvid-tester-apiscope"
#     description = "Scope opprettet fra test av Elvid Terraform provider (terraform-tester i terraform-provider-elvid)"
#     user_claims = ["email", "ad_groups"]
#     allow_user_clients = true
# }

## Module userclient
## Note that this require vault setup. Se readme
module "elvid_userclient" {
  source      = "C:\\3lvia\\terraform-elvid-userclient"
  environment = "dev"
  client_name = "test-bf"
  scopes = [ "louvre.imageapi.useraccess", "openid", "ad_groups"]
  domains = var.domains[var.environment]
  redirect_uri_paths = [ "/silentcallback.html", "/oidc/callback"]
  post_logout_redirect_uri_paths = [""]
  elvia_ad_login_enabled         = true
  system_name      = "elvid"
  client_secret_enabled = true
  ad-groups-filter = null
}

## Module machineclient
## Note that this require vault setup. Se readme
# module "elvid_machineclient" {
#   # source  = "app.terraform.io/Elvia/machineclient/elvid"
#   source      = "C:\\3lvia\\terraform-elvid-machineclient"
#   scopes = ["louvre.imageapi"]
#   environment      = var.environment
#   system_name      = "elvid"
#   application_name = "demo-machineclient2"
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
  type = map
  default = {
    "dev" = ["http://localhost:{port}", "https://examplesystem.dev-elvia.io"]
    "test" = ["https://examplesystem.test-elvia.io"]
    "prod" = ["https://examplesystem.elvia.io"]
  }
}
