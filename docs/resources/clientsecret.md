# ClientSecret Resource

This resource manages a clientsecret.

It is always used together with at machineclient or userclient.

## Example Usage
```hcl
resource "elvid_machineclient" "machineclient" {
   name = "example-machineclient"
   scopes = ["temp"]
}

resource "elvid_clientsecret" "clientsecret" {
    client_id = elvid_machineclient.machineclient.id
}
```

->Usage for this in Elvia is mostly done indirectly through a module.