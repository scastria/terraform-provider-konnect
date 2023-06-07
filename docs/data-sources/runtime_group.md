---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_runtime_group
Represents a runtime group
## Example usage
```hcl
data "konnect_runtime_group" "example" {
  name = "TestRuntimeGroup"
}
```
## Argument Reference
* `search_name` - **(Optional, String)** The search string to apply to the name of the runtime group. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the runtime group. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
* `description` - **(String)** The description of the runtime group.
* `cluster_type` - **(String)** The cluster type of the runtime group.
* `control_plane_endpoint` - **(String)** The control plane endpoint URL of the runtime group.
* `telemetry_endpoint` - **(String)** The telemetry endpoint URL of the runtime group.
