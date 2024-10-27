---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer_key
Represents an API key credential for a consumer within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_consumer" "Consumer" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_username = "Bob"
}
resource "konnect_consumer_key" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  consumer_id = data.konnect_consumer.Consumer.consumer_id
  key = "my-api-key"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `consumer_id` - **(Required, String)** The id of the consumer.
* `key` - **(Optional/Computed, String)** The API key value.  If left out, a key will be generated for you.
* `tags` - **(Optional, List of String)** An extra list of tags to assign to the API key in addition to the `default_tags` configured in the provider.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`consumer_id`:`key_id`
* `key_id` - **(String)** Id of the consumer API key alone
* `all_tags` - **(List of String)** The complete list of tags assigned to the API key, including the `tags` defined on this resource and the `default_tags` configured in the provider.
## Import
Consumer keys can be imported using a proper value of `id` as described above
