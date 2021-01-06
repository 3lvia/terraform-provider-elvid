# ElvID provider

This custom terraform provider is used to manage resources for ElvID, which is an Elvia application that uses IdentityServer.

It can be used to manage machineclients (client_credentials/password) and userclients (authorication code), along with ClientSecrets for these clients. It uses Azure AD service principal for authentication, and on the API side we requre a specific scope for authorization.

->Note that naming of resources, the authentication/authorization, and some schema variables (and their defaults) might be specific to Elvia's usecase. Still we hope you can use this for inspiration and as a base to create your custom provider to manage IdentityServer resources.

More info: https://github.com/3lvia/terraform-provider-elvid