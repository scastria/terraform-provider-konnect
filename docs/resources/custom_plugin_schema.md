---
subcategory: "Runtime Configuration"
---
# Resource: konnect_custom_plugin_schema
Represents a custom plugin schema within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
resource "konnect_custom_plugin_schema" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  name = "my-plugin"
  schema_lua = "return { name=\"my-plugin\", fields = { { config = { type = \"record\", fields = { } } } } }"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `name` - **(Required, ForceNew, String)** The name of the custom plugin schema.
* `schema_lua` - **(Required, String)** The lua code that defines the schema.  Typically, this is the content of the custom plugin schema.lua file.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`name`
## Import
Custom plugin schemas can be imported using a proper value of `id` as described above
