---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer_jwt
Represents a JWT credential for a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "Consumer" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
resource "konnect_consumer_jwt" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  consumer_id = data.konnect_consumer.Consumer.consumer_id
  key = "my-key"
  secret = "my-secret"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `consumer_id` - **(Required, String)** The id of the consumer.
* `key` - **(Optional/Computed, String)** The key value.  If left out, a key will be generated for you.
* `secret` - **(Optional/Computed, String)** The secret value.  If left out, a key will be generated for you.
* `algorithm` - **(Optional, String)** The algorithm for the JWT.  Allowed values: `HS256`, `HS384`, `HS512`, `RS256`, `RS384`, `RS512`, `ES256`, `ES384`. Default: `HS256`
* `rsa_public_key` - **(Optional, String)** The RSA public key in PEM format for the JWT.  Required if `algorithm` is `RS256`, `RS384`, `RS512`, `ES256`, or `ES384`.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`:`jwt_id`
* `jwt_id` - **(String)** Id of the consumer JWT alone
## Import
Consumer JWTs can be imported using a proper value of `id` as described above
