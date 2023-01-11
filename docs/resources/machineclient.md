# MachineClient Resource
 
This resource manages a machineclient (client_credentials openid-connect-client).

It is always used together with at least one clientsecret-resource.

## Example Usage
```hcl
resource "elvid_machineclient" "machineclient" {
   name = "example-machineclient"
   scopes = ["temp"]
   client_claims {
      type = "edna_topics_read"
      values = ["topicA"]
   }
   client_claims {
      type = "edna_topics_write"
      values = ["topicA", "topicB"]
   }
}
```

->Usage for this in Elvia is mostly done indirectly through a module.