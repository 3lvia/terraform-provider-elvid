# This is used to test the provider. Check the readme for info about installing and running this

# Provider
provider "elvid" {
  tenant_id                  = var.tenant_id
  terraform_sp_client_id     = var.terraform_sp_client_id
  terraform_sp_client_secret = var.terraform_sp_client_secret
  environment                = var.environment
}

## Machine client
resource "elvid_machineclient" "machineclient" {
  name   = "image-scanner"
  scopes = ["louvre.imageapi"]
}

resource "elvid_clientsecret" "clientsecret" {
  client_id = elvid_machineclient.machineclient.id
}

## Here we output the resources to console, normaly it would be distributed to vault
output "machineclient" {
  value = elvid_machineclient.machineclient
}

output "clientsecret" {
  value = elvid_clientsecret.clientsecret
}
