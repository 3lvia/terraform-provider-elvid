# This is used to test the provider. Check the readme for info about installing and running this

# Provider
provider "elvid" {
  tenant_id = var.tenant_id
  terraform_sp_client_id = var.terraform_sp_client_id
  terraform_sp_client_secret = var.terraform_sp_client_secret
  environment = var.environment
  # override_elvid_authority = "https://localhost:44383"
  override_elvid_authority = "https://elvid.dev-elvia.io"
  run_hashed_secret_validation = true
}

## User client

# resource "elvid_userclient" "userclient" {
#     client_name = "jordfeil-ui"
#     scopes = ["temp", "openid", "ad_groups"]
#     domains = var.domains[var.environment]
#     redirect_uri_paths = ["/callback.html"]
#     post_logout_redirect_uri_paths = ["/index.htm"]
#     bankid_login_enabled = true
#     local_login_enabled = true
#     elvia_ad_login_enabled = false
#     test_user_login_enabled = false
#     require_client_secret = false
#     access_token_life_time = 3598
#     always_include_user_claims_in_id_token = true
#     client_name_language_key = null
#     allow_use_of_refresh_tokens = false
#     one_time_usage_for_refresh_tokens = true
#     refresh_token_life_time = 2592000
# }

## Machine client

# resource "elvid_machineclient" "machineclient" {
#     name = "onsdag"
#     test_user_login_enabled = true
#     access_token_life_time = 3517
#     resource_taint_version = "1"
#     scopes = ["temp"]
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
#     allow_machine_clients = true
# }
