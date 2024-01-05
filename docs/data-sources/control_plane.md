---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_control_plane
Represents a control plane
## Example usage
```hcl
data "konnect_control_plane" "example" {
  name = "TestControlPlane"
}
```
## Argument Reference
* `search_name` - **(Optional, String)** The search string to apply to the name of the control plane. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the control plane. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
* `description` - **(String)** The description of the control plane.
* `cluster_type` - **(String)** The cluster type of the control plane.
* `control_plane_endpoint` - **(String)** The control plane endpoint URL of the control plane.
* `telemetry_endpoint` - **(String)** The telemetry endpoint URL of the control plane.
