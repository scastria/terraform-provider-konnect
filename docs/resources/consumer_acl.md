---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer_acl
Represents an ACL credential for a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "Consumer" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
resource "konnect_consumer_acl" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  consumer_id = data.konnect_consumer.Consumer.consumer_id
  group = "my-acl-group"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `consumer_id` - **(Required, String)** The id of the consumer.
* `group` - **(Required, String)** The ACL group value.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`:`acl_id`
* `acl_id` - **(String)** Id of the consumer ACL alone
## Import
Consumer ACLs can be imported using a proper value of `id` as described above
