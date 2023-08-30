---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_consumer
Represents a consumer
## Example usage
```hcl
data "konnect_runtime_group" "RuntimeGroup" {
  name = "TestRuntimeGroup"
}
data "konnect_consumer" "example" {
  runtime_group_id = data.konnect_runtime_group.RuntimeGroup.id
  search_username = "Bob"
}
```
## Argument Reference
* `runtime_group_id` - **(Required, String)** The id of the runtime group containing consumer
* `search_username` - **(Optional, String)** The search string to apply to the username of the consumer. Uses contains.
* `username` - **(Optional, String)** The filter string to apply to the username of the consumer. Uses equality.
* `search_custom_id` - **(Optional, String)** The search string to apply to the custom_id of the consumer. Uses contains.
* `custom_id` - **(Optional, String)** The filter string to apply to the custom_id of the consumer. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `runtime_group_id`:`consumer_id`
* `consumer_id` - **(String)** Id of the consumer alone
