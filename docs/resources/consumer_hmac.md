---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer_hmac
Represents an HMAC credential for a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "Consumer" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
resource "konnect_consumer_hmac" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  consumer_id = data.konnect_consumer.Consumer.consumer_id
  username = "my-username"
  secret = "my-secret"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `consumer_id` - **(Required, String)** The id of the consumer.
* `username` - **(Required, String)** The username value.
* `secret` - **(Optional/Computed, String)** The secret value.  If left out, a secret will be generated for you.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`:`hmac_id`
* `hmac_id` - **(String)** Id of the consumer HMAC alone
## Import
Consumer HMACs can be imported using a proper value of `id` as described above
