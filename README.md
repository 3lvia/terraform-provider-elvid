# terraform-provider-elvid

This custom terraform provider is used to manage resources for ElvID, which is an Elvia application that uses IdentityServer.

The provider is published to [registry.terraform.io/providers/3lvia/elvid](https://registry.terraform.io/providers/3lvia/elvid/latest).

It can be used to manage machineclients (client_credentials/password), userclients (authorication code), ClientSecrets for these clients and API scopes.
It uses Azure AD service principal for authentication, and on the API side we require a specific scope for authorization.

Note that naming of resources, the authentication/authorization, and some schema variables (and their defaults) might be  specific to Elvia's usecase.
Still we hope you can use this repo for inspiration and as a base to create your custom provider to manage IdentityServer resources.

To implement this in your IdentityServer solution you also need to create the API that receives these requests and updates your ConfigurationStore.

# General information about creating a custom terraform provider

Se [here](https://learn.hashicorp.com/collections/terraform/providers) for the general information about creating custom providers from Terraform.

# Local Setup
## Install go
[Golang installation guide](https://golang.org/doc/install)

## Install terraform
[Install terraform](https://learn.hashicorp.com/terraform/getting-started/install.html). For windows you can add terraform.exe to {user}/bin. Make sure %USERPROFILE%\go\bin is in path, and above the go-spesific paths.

## Checkout code
Checkout the code-repo to {GOPATH}\src\github.com\3lvia\terraform-provider-elvid

# Project structure
* repo-root
  * terraform-tester.tf/versions.tf: terraform files for manually testing the provider.
  * elvidapiclient: go class-library for getting AccessToken from AD and calling ElvID-api
  * main.go: Standard file, sets up serving of the provider by calling the Provider()-function.
  * provider.go: Defines the provider schema (inputs to the provider), the mapping to resorces, and the interface that is passed to resrouces
  * resource_clientsecret.go: Defines the resource schema and methods for clientsecrets.
  * resource_machineclient.go: Defines the resource schema and methods for machineclient.
  * resource_userclient.go: Defines the resource schema and methods for userclient.
  * resource_apiscope.go: Defines the resource schema and methods for apiscope.

# Setup terraform for running locally
## Setup dev overrides to target local build of the provider
This makes sure that terraform will use the local build of the provider and not the published build from terraform registry.

Note that terraform init will download the published library from terraform registry, but the dev_overrides variant will still be used on plan/apply. 

This is done in the .terraformrc/terraform.rc file see https://www.terraform.io/cli/config/config-file#locations

For windows create/edit $env:APPDATA\terraform.rc and add provider_installation

```console
provider_installation {
  # Override provider for local development
  dev_overrides {
    "3lvia/elvid" = "C:\\Users\\{{replace with your windows username}}\\go\\src\\github.com\\3lvia\\terraform-provider-elvid"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

## Adding terraform.tfvars to terraform-tester
Create {repo-root}/terraform.tfvars and add these variables.
Secret values can be found in vault-dev in the path /elvid/kv/manual/terraform_provider_elvid_adcredentials

```
terraform_sp_client_id = "replaceme"
terraform_sp_client_secret = "replaceme"
tenant_id = "replaceme"
```

Note that terraform.tfvars is added to .gitignore. Make sure to newer publish these secrets. This is a public repository.

# Running locally

## Build the provider for a local run

```console
# from repo-root
go build
```

## Running terraform locally
Make sure you have setup terraform for running locally (described above)

```console
# from repo-root
terraform apply;
```

You should get a warning on plan / apply
```console
 The following provider development overrides are set in the CLI configuration:
│  - 3lvia/elvid in C:\Users\{{username}}\go\src\github.com\3lvia\terraform-provider-elvid
```

You don't usually need to run terraform init because we are using dev_overrides.
If you are working with modules, you might have to do terraform init (it will tell you when running plan or apply).
Terraform init will download the published library from terraform registry, but the dev_overrides variant will still be used on plan/apply. 

## Build and apply with one command
```console
# from repo-root
go build; terraform apply -auto-approve;
```

# Debugging
Debugging the go-code when running from terraform is not added to this repo. See [this guide if debugging should be considered](https://developer.hashicorp.com/terraform/plugin/sdkv2/guides/v2-upgrade-guide#support-for-debuggable-provider-binaries).

It is possible to print debug info as warnings in diag.Diagnostics. This is used for ApiScope. It requires v2 of the SDK, and some rewrite of the resource definition, as in resource_apiscope.go/apiscopeservice.go. See [the upgrade guide for v2 of the SDK](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html). Terraform-privider-elvid already uses v2, but v2 also supports the v1 way.

For resources/services that is not yet rewritten to v2 (but still use error and Create instead of CreateContext), debugging can be done by writing a file with debug messages:

```
# for a string 
ioutil.WriteFile("custom-log.text", []byte(someString), 0644)

# for an object
serialized, _ := json.Marshal(someObject)
ioutil.WriteFile("custom-log.text", []byte(serialized), 0644)
```
# Publish a new release
## Publish to terraform registry
To publish to [registry.terraform.io/providers/3lvia/elvid](https://registry.terraform.io/providers/3lvia/elvid/latest) create a new github-release in this repo. 
Github-actions is setup to atomatically build and publish new releases. 

Github-actions uses the our private signing key to sign the build. The public variant of this key is added in terraform registry.
Backup of the key is found in vault (prod) elvid/kv/manual/elvid-provider-build-signing-key

# Notes from creating a custom provider
* Creating a class-library to wrap the api was helpful, to get more clean resource-code. It was benefitial to have it in the same repo. 
* The code must be in {GOPATH}\src\github.com\3lvia
* Filename resources must be in format resource_{resourcename}.go
* Creating a "resource_taint_version" variable with ForceNew=true was very helpfull to quickly test changes, and will be helpfull in actual use as well. 
* Terraform does handling state and knows when to call create, read, delete, update. So create good variable-schemas, implement these methods and let terraform handle the rest.
* The id field has to be Optional and Computed, so even resources where the id can be defined in the tf file, it will be "(known after apply)". Example: apiscope, where id=name, and we use the name field as required input, and set id=name when the resource is created.
