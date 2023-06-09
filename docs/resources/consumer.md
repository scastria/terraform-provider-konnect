---
subcategory: "Runtime Configuration"
---
# Resource: konnect_consumer
Represents a consumer within a runtime group
## Example usage
```hcl
data "konnect_runtime_group" "RuntimeGroup" {
  name = "TestRuntimeGroup"
}
resource "konnect_consumer" "example" {
  runtime_group_id = data.konnect_runtime_group.RuntimeGroup.id
  username = "testuser"
  custom_id = "testcustom"
}
```
## Argument Reference
* `runtime_group_id` - **(Required, String)** The id of the runtime group.
* `username` - **(Optional, String)** The unique username of the Consumer.
* `custom_id` - **(Optional, String)** Field for storing an existing unique ID for the Consumer.
## Attribute Reference
* `id` - **(String)** Same as `runtime_group_id`:`consumer_id`
* `consumer_id` - **(String)** Id of the consumer alone
## Import
Consumers can be imported using a proper value of `id` as described above
