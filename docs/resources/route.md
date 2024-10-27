---
subcategory: "Runtime Configuration"
---
# Resource: konnect_route
Represents a route within a control plane
## Example usage
```hcl
data "konnect_control_plane" "ControlPlane" {
  name = "TestControlPlane"
}
resource "konnect_service" "Service" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  host = "mockbin.org"
  name = "Test"
}
resource "konnect_route" "example" {
  control_plane_id = data.konnect_control_plane.ControlPlane.id
  service_id = konnect_service.Service.service_id
  name = "Test"
  protocols = ["http"]
  paths = ["/example"]
  header {
    name = "required-header"
    values = ["required-header-values"]
  }
}
```
## Argument Reference
* `control_plane_id` - **(Required, String)** The id of the control plane.
* `service_id` - **(Optional, String)** The id of the service to forward traffic to.
* `name` - **(Optional, String)** The name of the route.
* `protocols` - **(Optional, List of String)** The protocols this route should allow. Allowed values: `http`, `https`
* `methods` - **(Optional, List of String)** The methods this route should allow. Allowed values: `GET`, `PUT`, `POST`, `PATCH`, `DELETE`, `OPTIONS`, `HEAD`, `CONNECT`, `TRACE`
* `hosts` - **(Optional, List of String)** The hosts this route should allow.
* `paths` - **(Optional, List of String)** The paths this route should allow.
* `https_redirect_status_code` - **(Optional, Integer)** The status code Kong responds with when all properties of a Route match except the protocol. Allowed values: `426`, `301`, `302`, `307`, `308`. Default: `426`
* `regex_priority` - **(Optional, Integer)** A number used to choose which route resolves a given request when several routes match it using regexes simultaneously. Default: `0`
* `strip_path` - **(Optional, Boolean)** Whether to strip the matching prefix from the Service request. Default: `true`
* `path_handling` - **(Optional, String)** Controls how the Service path, Route path and requested path are combined when sending a request to the upstream. Allowed values: `v0`, `v1`. Default: `v0`
* `preserve_host` - **(Optional, Boolean)** Whether to use the request `Host` header during the Service request. Default: `false`
* `request_buffering` - **(Optional, Boolean)** Whether to enable request body buffering. Default: `true`
* `response_buffering` - **(Optional, Boolean)** Whether to enable response body buffering. Default: `true`
* `header` - **(Optional, set{header})** Configuration block for a header.  Can be specified multiple times for each header.  Each block supports the fields documented below.
* `tags` - **(Optional, List of String)** An extra list of tags to assign to the route in addition to the `default_tags` configured in the provider.
### header
* `name` - **(Required, String)** Name of header this route should require.
* `values` - **(Required, List of String)** Allowed values this header should equal.
## Attribute Reference
* `id` - **(String)** Same as `control_plane_id`:`route_id`
* `route_id` - **(String)** Id of the route alone
* `all_tags` - **(List of String)** The complete list of tags assigned to the route, including the `tags` defined on this resource and the `default_tags` configured in the provider.
## Import
Routes can be imported using a proper value of `id` as described above
