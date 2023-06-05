# Resource: konnect_runtime_group
Represents a runtime group
## Example usage
```hcl
resource "konnect_runtime_group" "example" {
  name = "TestRuntimeGroup"
  description = "TestRuntimeGroup"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the runtime group.
* `description` - **(Optional, String)** The description of the runtime group.
## Attribute Reference
* `id` - **(String)** Guid
* `cluster_type` - **(String)** The cluster type of the runtime group.
* `control_plane_endpoint` - **(String)** The control plane endpoint URL of the runtime group.
* `telemetry_endpoint` - **(String)** The telemetry endpoint URL of the runtime group.
## Import
Runtime groups can be imported using a proper value of `id` as described above
