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
  * terraform-tester: folder with terraform files for manually testing the provider.
  * elvidapiclient: go class-library for getting AccessToken from AD and calling ElvID-api
  * main.go: Standard file, sets up serving of the provider by calling the Provider()-function.
  * provider.go: Defines the provider schema (inputs to the provider), the mapping to resorces, and the interface that is passed to resrouces
  * resource_clientsecret.go: Defines the resource schema and methods for clientsecrets.
  * resource_machineclient.go: Defines the resource schema and methods for machineclient.
  * resource_userclient.go: Defines the resource schema and methods for userclient.
  * resource_apiscope.go: Defines the resource schema and methods for apiscope.

# Running locally

## Build the provider for a local run
The provider must be installed by building it to one of the [plugin locations that terraform init searches through](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations).

Note that "terraform init" searches the current directory. So one could have everything in root here. 
That is nice to get things started, but it got a bit messy, so I moved the terraform files to a seperate folde (terraform-tester).

So instead install the provider in one of the common provicer locations.  
```console
# Windows (from repo-root)
go build -o %APPDATA%\terraform.d\plugins\local\3lvia\elvid\9999.9.9\windows_amd64\terraform-provider-elvid_v9999.9.9.exe
or
go build -o $env:APPDATA\terraform.d\plugins\local\3lvia\elvid\9999.9.9\windows_amd64\terraform-provider-elvid_v9999.9.9.exe

# Linux (from repo-root)
go build -o ~/.terraform.d/plugins/local/3lvia/elvid/9999.9.9/linux_amd64/terraform-provider-elvid_v9999.9.9
```

## Adding terraform.tfvars to terraform-tester
Create ../terraform-tester/terraform.tfvars and add these variables.
Secret values can be found in vault-dev in the path /elvid/kv/manual/terraform_provider_elvid_adcredentials

```
terraform_sp_client_id = "replaceme"
terraform_sp_client_secret = "replaceme"
tenant_id = "replaceme"
```

Note that terraform.tfvars is added to .gitignore. Make sure to newer publish these secrets. This is a public repository.

## Running terraform locally

```console
# from repo-root/terraform-tester
terraform init;
terraform apply;
```
Note: terraform will not fetch the new version of the go library when you build again, when using "v9999.9.9" repeatedly. You need to delete the .terraform.lock.hcl file in your terraform-tester folder and run terraform init again to run with the newly built version (the .terraform directory contains a copy of the previous version of the go library, but you don't need to delete it to get your new build).

## Build and run in one command
```console
# from repo-root/terraform-tester
rm .terraform.lock.hcl -ErrorAction Ignore; cd ..; go build -o $env:APPDATA\terraform.d\plugins\local\3lvia\elvid\9999.9.9\windows_amd64\terraform-provider-elvid_v9999.9.9.exe; cd .\terraform-tester; terraform init; terraform apply
```
See the note in "Running terraform locally" about deleting the terraform lock file after building new "v9999.9.9" version.

# Debugging
Debugging the go-code when running from terraform is not suported. It is possible to print debug info as warnings in diag.Diagnostics. This is used for ApiScope. It requires v2 of the SDK, and some rewrite of the resource definition, as in resource_apiscope.go/apiscopeservice.go. See [the upgrade guide for v2 of the SDK](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html). Terraform-privider-elvid already uses v2, but v2 also supports the v1 way.

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

## Publish to terraform-plugins (The old method without terraform registry for terraform <= 0.12)
For terraform 12 we can't read from terraform registry, instead we have a repo (terraform-plugins) where we publish the compiled binaries.
The terraform-plugins repo is added to the terreform-repos using git submodule.

In windows powershell (or adust for linux)
```console
$env:BuildProviderVersion = '1.0.3'
$env:GOOS = "windows";$env:GOARCH = "amd64"; go build -o C:\3lvia\terraform-plugins\windows_amd64\terraform-provider-elvid_v$env:BuildProviderVersion.exe
$env:GOOS = "linux"; $env:GOARCH = "amd64"; go build -o C:\3lvia\terraform-plugins\linux_amd64\terraform-provider-elvid_v$env:BuildProviderVersion
$env:GOOS = "darwin";$env:GOARCH = "amd64"; go build -o C:\3lvia\terraform-plugins\darwin_amd64\terraform-provider-elvid_v$env:BuildProviderVersion
Copy-Item "C:\3lvia\terraform-plugins\linux_amd64\terraform-provider-elvid_v$env:BuildProviderVersion" -Destination "C:\3lvia\terraform-plugins\terraform-provider-elvid_v$env:BuildProviderVersion"
```

Follow the README in terraform-plugins to finish publishing there.

# Notes from creating a custom provider
* Creating a class-library to wrap the api was helpful, to get more clean resource-code. It was benefitial to have it in the same repo. 
* The code must be in {GOPATH}\src\github.com\3lvia
* Filename resources must be in format resource_{resourcename}.go
* Creating a "resource_taint_version" variable with ForceNew=true was very helpfull to quickly test changes, and will be helpfull in actual use as well. 
* Terraform does handling state and knows when to call create, read, delete, update. So create good variable-schemas, implement these methods and let terraform handle the rest.
* The id field has to be Optional and Computed, so even resources where the id can be defined in the tf file, it will be "(known after apply)". Example: apiscope, where id=name, and we use the name field as required input, and set id=name when the resource is created.
