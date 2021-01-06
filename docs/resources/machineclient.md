# MachineClient Resource
 
This resource manages a machineclient (client_credentials openid-connect-client).

It is always used together with at least one clientsecret-resource.

## Example Usage
```hcl
resource "elvid_machineclient" "machineclient" {
   name = "example-machineclient"
   scopes = ["temp"]
}
```

->Usage for this in Elvia is mostly done indirectly through a module.