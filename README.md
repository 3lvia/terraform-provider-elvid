# terraform-provider-elvid

This custom terraform provider is used to manage resources for ElvID, which is an Elvia application that uses IdentityServer.

It can be used to manage machineclients (client_credentials/password) and userclients (authorication code), along with ClientSecrets for these clients.
It uses Azure AD service principal for authentication, and on the API side we requre a specific scope for authorization.

Note that naming of resources, the authentication/authorization, and some schema variables (and their defaults) might be  specific to Elvia's usecase.
Still we hope you can use this repo for inspiration and as a base to create your custom provider to manage IdentityServer resources.

To implement this in your IdentityServer solution you also need to create the API that receives these requests and updates your ConfigurationStore.

# General information about creating a custom terraform provider

Se [here](https://learn.hashicorp.com/collections/terraform/providers) for the general information about creating custom providers from Terraform 

# Local Setup
## Install go
[Golang installation guide](https://golang.org/doc/install)

## Install terraform
[Install terraform](https://learn.hashicorp.com/terraform/getting-started/install.html).For windows you can add terraform.exe to {user}/bin.Make sure %USERPROFILE%\go\bin is in path, and above the go-spesific paths.

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

# Running locally

## Buld the provider for a local run
The provider must be installed by building it to one of the [plugin locations that terraform init searches through](https://www.terraform.io/docs/extend/how-terraform-works.html#plugin-locations).

Note that "terraform init" searches the current directory. So one could have everything in root here. 
That is nice to get things started, but it got a bit messy, so I moved the terraform files to a seperate folde (terraform-tester).

So instead install the provider in one of the common provicer locations.  
```console
# Windows (from repo-root)
go build -o $env:APPDATA\terraform.d\plugins\terraform-provider-elvid.exe

# Linux (from repo-root)
go build -o ~/.terraform.d/plugins/terraform-provider-elvid
```

## Adding terraform.tfvars to terraform-tester
Create ../terraform-tester/terraform.tfvars and add these variables.
Secret values can be found in vault-dev in the path /elvid/kv/manual/terraform_provider_elvid_adcredentials

```
terraform_sp_client_id = "replaceme"
terraform_sp_client_secret = "replaceme"
tenant_id = "replaceme"
```

## Running terraform locally

```console
# from repo-root/terraform-tester
terraform init;
terraform apply;
```

## Build and run in one command
```console
# from repo-root/terraform-tester
cd ..; go build -o $env:APPDATA\terraform.d\plugins\terraform-provider-elvid.exe; cd .\terraform-tester; terraform init; terraform apply
```

# Debugging
Debugging the go-code when running from terraform is not suported. 
The best way I found was to write a file with the message I needed to check (a suggestion for a better solution here would be very helpfull)

```
# for a string 
ioutil.WriteFile("custom-log.text", []byte(someString), 0644)

# for a object
serialized, _ := json.Marshal(someObject)
ioutil.WriteFile("custom-log.text", []byte(serialized), 0644)
```
# Publish a new release
## Publish to terraform registry
Work in progress

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
* Creating a class-library to wrap the api was helpful, to get more clean resource-code. I was benefitial to have it in the same repo. 
* The code must be in {GOPATH}\src\github.com\3lvia
* Filename resources must be in format resource_{resourcename}.go
* Creating a "resource_taint_version" variable with ForceNew=true was very helpfull to quickly test changes, and will be helpfull in actual use as well. 
* Terraform does handling state and knows when to call create, read, delete, update. So create good variable-schemas, implement these methods and let terraform handle the rest