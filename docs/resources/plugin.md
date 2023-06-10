---
subcategory: "Runtime Configuration"
---
# Resource: konnect_plugin
Represents a plugin within a runtime group
## Example usage
```hcl
data "konnect_runtime_group" "RuntimeGroup" {
  name = "TestRuntimeGroup"
}
resource "konnect_plugin" "example" {
  runtime_group_id = data.konnect_runtime_group.RuntimeGroup.id
  name = "rate-limiting"
  config_json = <<EOF
{
  "minute": 5
}
EOF
}
```
## Argument Reference
* `runtime_group_id` - **(Required, String)** The id of the runtime group.
* `name` - **(Required, String)** The name of the plugin which must match a valid installed plugin.
* `instance_name` - **(Optional, String)** The instance name of the plugin. Default: `-`
* `protocols` - **(Optional, List of String)** A list of the request protocols that will trigger this plugin. Allowed values: `grpc`, `grpcs`, `http`, `https`, `tcp`, `tls`, `tls_passthrough`, `udp`, `ws`, `wss`
* `enabled` - **(Optional, Boolean)** Whether the plugin is active. Default: `true`
* `config_json` - **(Optional, JSON)** The configuration properties for the plugin which can be found on the plugins documentation page in the [Kong Plugin Hub](https://docs.konghq.com/hub/) 
* `service_id` - **(Optional, String)** If set, the plugin will only activate when receiving requests via one of the routes belonging to the specified service.
* `route_id` - **(Optional, String)** If set, the plugin will only activate when receiving requests via the specified route.
* `consumer_id` - **(Optional, String)** If set, the plugin will activate only for requests where the specified consumer has been authenticated.
## Attribute Reference
* `id` - **(String)** Same as `runtime_group_id`:`plugin_id`
* `plugin_id` - **(String)** Id of the plugin alone
* `config_all_json` - **(JSON)** The full configuration properties for the plugin, including all properties with their default values not specified in `config`.
## Import
Plugins can be imported using a proper value of `id` as described above
