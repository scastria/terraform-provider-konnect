---
subcategory: "Runtime Configuration"
---
# Resource: konnect_service
Represents a service within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
resource "konnect_service" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  host = "mockbin.org"
  name = "Test"
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `host` - **(Required, String)** The host of the service.
* `name` - **(Optional, String)** The name of the service.
* `retries` - **(Optional, Integer)** The number of retries to execute upon failure to proxy. Default: `5`
* `protocol` - **(Optional, String)** The protocol used to communicate with the host. Allowed values: `grpc`, `grpcs`, `http`, `https`, `tcp`, `tls`, `tls_passthrough`, `udp`, `ws`, `wss`
* `port` - **(Optional, Integer)** The port used to communicate with the host. Default: `80`
* `path` - **(Optional, String)** The path to be used in requests to the host.
* `connect_timeout` - **(Optional, Integer)** The timeout in milliseconds for establishing a connection to the host. Default: `60000`
* `read_timeout` - **(Optional, Integer)** The timeout in milliseconds between two successive read operations for transmitting a request to the host. Default: `60000`
* `write_timeout` - **(Optional, Integer)** The timeout in milliseconds between two successive write operations for transmitting a request to the host. Default: `60000`
* `enabled` - **(Optional, Boolean)** Whether the service is active. Default: `true`
* `tags` - **(Optional, List of String)** An extra list of tags to assign to the service in addition to the `default_tags` configured in the provider.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`service_id`
* `service_id` - **(String)** Id of the service alone
* `all_tags` - **(List of String)** The complete list of tags assigned to the service, including the `tags` defined on this resource and the `default_tags` configured in the provider.
## Import
Services can be imported using a proper value of `id` as described above
