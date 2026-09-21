# ElvID provider

This custom terraform provider is used to manage resources for ElvID, which is an Elvia application that uses IdentityServer.

It can be used to manage machineclients (client_credentials/password) and userclients (authorication code), along with ClientSecrets for these clients. It uses Azure AD service principal for authentication, and on the API side we require a specific scope for authorization.

->Note that naming of resources, the authentication/authorization, and some schema variables (and their defaults) might be specific to Elvia's usecase. Still we hope you can use this for inspiration and as a base to create your custom provider to manage IdentityServer resources.

More info: https://github.com/3lvia/terraform-provider-elvid

## Authentication

The provider logs in to Entra as a service principal and exchanges that for an access token for the ElvID API. Two ways, the first preferred:

- **Federated token (no secret).** The app registration trusts the run's platform (Scalr, Terraform Cloud, GitHub Actions) as a federated identity credential, and the provider sends the platform's OIDC token as `client_assertion`. Set `terraform_sp_client_assertion`, or leave it empty and let the provider read `ARM_OIDC_TOKEN` / `ARM_OIDC_TOKEN_FILE_PATH`, the same variables the azurerm provider uses, so a Scalr workspace with an Azure OIDC provider configuration needs nothing extra.
- **Client secret.** Set `terraform_sp_client_secret`. Legacy; the secret has to be stored and rotated somewhere.

```terraform
provider "elvid" {
  tenant_id              = "..."
  terraform_sp_client_id = var.client_id
  # terraform_sp_client_secret left empty: the run's federated token (ARM_OIDC_TOKEN) is used.
  environment = var.environment
}
```

