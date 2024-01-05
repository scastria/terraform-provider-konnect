---
subcategory: "Runtime Configuration"
---
# Resource: konnect_control_plane
Represents a control plane
## Example usage
```hcl
resource "konnect_control_plane" "example" {
  name = "TestControlPlane"
  description = "TestControlPlane"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the control plane.
* `description` - **(Optional, String)** The description of the control plane.
## Attribute Reference
* `id` - **(String)** Guid
* `cluster_type` - **(String)** The cluster type of the control plane.
* `control_plane_endpoint` - **(String)** The control plane endpoint URL of the control plane.
* `telemetry_endpoint` - **(String)** The telemetry endpoint URL of the control plane.
## Import
Control planes can be imported using a proper value of `id` as described above
