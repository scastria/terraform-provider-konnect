---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_nodes
Represents all nodes of a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_nodes" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the parent control plane.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`
* `nodes` - **(set{node})** Set of nodes belonging to control plane
### node
* `id` - **(String)** Id of node.
* `version` - **(String)** Version of node.
* `hostname` - **(String)** Hostname of node.
* `last_ping` - **(Integer)** Last time of ping of node.
* `type` - **(String)** Type of node.
* `config_hash` - **(String)** Hash of the current configuration state of node.
* `data_plane_cert_id` - **(String)** Id of certificate used in communication between node and control plane.
