# Resource: konnect_identity_provider
Represents Konnect identity provider settings
## Example usage
```hcl
resource "konnect_identity_provider" "example" {
  issuer = "https://example.com"
  client_id = "XXXX"
  login_path = "login"
  scopes = [
    "email",
    "openid",
    "profile"
  ]
  email_claim_mapping = "email"
  name_claim_mapping = "name"
  groups_claim_mapping = "groups"
}
```
## Argument Reference
* `issuer` - **(Optional, String)** Issuer of the identity provider.
* `client_id` - **(Optional, String)** Client ID of the identity provider.
* `login_path` - **(Optional, String)** Login path of the identity provider.
* `scopes` - **(Optional, List of String)** Scopes of the identity provider.
* `email_claim_mapping` - **(Optional, String)** Claim to map email for the identity provider.
* `name_claim_mapping` - **(Optional, String)** Claim to map name for the identity provider.
* `groups_claim_mapping` - **(Optional, String)** Claim to map groups for the identity provider.
## Attribute Reference
* `id` - **(String)** Always equal to `identity-provider`
## Import
Identity provider can be imported using a proper value of `id` as described above
