---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_nodes
Represents all nodes of a runtime group
## Example usage
```hcl
data "konnect_runtime_group" "RuntimeGroup" {
  name = "TestRuntimeGroup"
}
data "konnect_nodes" "example" {
  runtime_group_id = data.konnect_runtime_group.RuntimeGroup.id
}
```
## Argument Reference
* `runtime_group_id` - **(Required, String)** The id of the parent runtime group.
## Attribute Reference
* `id` - **(String)** Same as `runtime_group_id`
* `nodes` - **(set{node})** Set of nodes belonging to runtime group
### node
* `id` - **(String)** Id of node.
* `version` - **(String)** Version of node.
* `hostname` - **(String)** Hostname of node.
* `last_ping` - **(Integer)** Last time of ping of node.
* `type` - **(String)** Type of node.
* `config_hash` - **(String)** Hash of the current configuration state of node.
* `data_plane_cert_id` - **(String)** Id of certificate used in communication between node and runtime group.
