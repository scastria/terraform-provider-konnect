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
* `username` - **(Optional, String)** The unique username of the consumer.
* `custom_id` - **(Optional, String)** Field for storing an existing unique ID for the consumer.
* `tags` - **(Optional, List of String)** An extra list of tags to assign to the consumer in addition to the `default_tags` configured in the provider.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`
* `consumer_id` - **(String)** Id of the consumer alone
* `all_tags` - **(List of String)** The complete list of tags assigned to the consumer, including the `tags` defined on this resource and the `default_tags` configured in the provider.
## Import
Consumers can be imported using a proper value of `id` as described above
