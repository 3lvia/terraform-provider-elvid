# UserClient Resource

This resource manages a useclient (Authorization code openid-connect-client).

## Example Usage
```hcl
resource "elvid_userclient" "userclient" {
   client_name = "example-userclient"
   scopes = ["temp"]
   domains = ["http://localhost:{port}", "https://examplesystem.dev-elvia.io"]
   redirect_uri_paths = ["/callback.html"]
   post_logout_redirect_uri_paths = ["/index.htm"]
   bankid_login_enabled = true
   local_login_enabled = true
   elvia_ad_login_enabled = false
   hafslund_ad_login_enabled = false
   test_user_login_enabled = false
   require_client_secret = false
   always_include_user_claims_in_id_token = true
   client_name_language_key = null
   allow_use_of_refresh_tokens = false
   one_time_usage_for_refresh_tokens = true
}
```

->Usage for this in Elvia is mostly done indirectly through a module.