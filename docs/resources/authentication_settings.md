---
subcategory: "Identity Management"
---
# Resource: konnect_authentication_settings
Represents authentication settings
## Example usage
```hcl
resource "konnect_authentication_settings" "example" {
  basic_auth_enabled = true
  oidc_auth_enabled = true
  idp_mapping_enabled = false
  konnect_mapping_enabled = true
}
```
## Argument Reference
* `basic_auth_enabled` - **(Optional, Boolean)** Whether basic authentication is enabled.
* `oidc_auth_enabled` - **(Optional, Boolean)** Whether OIDC authentication is enabled.
* `idp_mapping_enabled` - **(Optional, Boolean)** Whether IDP mapping is enabled.
* `konnect_mapping_enabled` - **(Optional, Boolean)** Whether Konnect mapping is enabled.
## Attribute Reference
* `id` - **(String)** Always equal to `authentication-settings`
## Import
Authentication settings can be imported using a proper value of `id` as described above
