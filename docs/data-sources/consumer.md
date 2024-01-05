---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_consumer
Represents a consumer
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane containing consumer
* `search_username` - **(Optional, String)** The search string to apply to the username of the consumer. Uses contains.
* `username` - **(Optional, String)** The filter string to apply to the username of the consumer. Uses equality.
* `search_custom_id` - **(Optional, String)** The search string to apply to the custom_id of the consumer. Uses contains.
* `custom_id` - **(Optional, String)** The filter string to apply to the custom_id of the consumer. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`
* `consumer_id` - **(String)** Id of the consumer alone
