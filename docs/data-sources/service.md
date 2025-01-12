---
subcategory: "Runtime Configuration"
---
# Data Source: konnect_service
Represents a service within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
data "konnect_service" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  search_name = "api"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `search_name` - **(Optional, String)** The search string to apply to the name of the service. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the service. Uses equality.
## Attribute Reference
* `retries` - **(Integer)** The number of retries to execute upon failure to proxy.
* `protocol` - **(Integer)** The protocol used to communicate with the host.
* `port` - **(Integer)** The port used to communicate with the host.
* `path` - **(String)** The path to be used in requests to the host.
* `connect_timeout` - **(Integer)** The timeout in milliseconds for establishing a connection to the host.
* `read_timeout` - **(Integer)** The timeout in milliseconds between two successive read operations for transmitting a request to the host.
* `write_timeout` - **(Integer)** The timeout in milliseconds between two successive write operations for transmitting a request to the host.
* `enabled` - **(Boolean)** Whether the service is active.
* `tags` - **(List of String)** An extra list of tags to assign to the service in addition to the `default_tags` configured in the provider.
* `id` - **(String)** Same as `control_plane_id`:`service_id`
* `service_id` - **(String)** Id of the service alone
* `all_tags` - **(List of String)** The complete list of tags assigned to the service, including the `tags` defined on this resource and the `default_tags` configured in the provider.
