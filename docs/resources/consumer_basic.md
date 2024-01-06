---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer_basic
Represents a basic auth credential for a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "Consumer" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
resource "konnect_consumer_basic" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  consumer_id = data.konnect_consumer.Consumer.consumer_id
  username = "my-username"
  password = "my-password"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `consumer_id` - **(Required, String)** The id of the consumer.
* `username` - **(Required, String)** The username value.
* `password` - **(Required, String)** The password value.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`:`basic_id`
* `password_hash` - **(String)** Hash of the password
* `basic_id` - **(String)** Id of the consumer basic auth alone
## Import
Consumer basics can be imported using a proper value of `id` as described above
