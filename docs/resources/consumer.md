---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer
Represents a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
resource "konnect_consumer" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  username = "testuser"
  custom_id = "testcustom"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `username` - **(Optional, String)** The unique username of the Consumer.
* `custom_id` - **(Optional, String)** Field for storing an existing unique ID for the Consumer.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`
* `consumer_id` - **(String)** Id of the consumer alone
## Import
Consumers can be imported using a proper value of `id` as described above
