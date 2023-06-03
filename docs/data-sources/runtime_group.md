# Data Source: konnect_runtime_group
Represents a Konnect runtime group
## Example usage
```hcl
data "konnect_runtime_group" "example" {
  name = "TestRuntimeGroup"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the runtime group.
## Attribute Reference
* `id` - **(String)** Guid
* `description` - **(String)** The description of the runtime group.
* `cluster_type` - **(String)** The cluster type of the runtime group.
* `control_plane_endpoint` - **(String)** The control plane endpoint URL of the runtime group.
* `telemetry_endpoint` - **(String)** The telemetry endpoint URL of the runtime group.
