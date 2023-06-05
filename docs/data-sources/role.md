# Data Source: konnect_role
Represents a role
## Example usage
```hcl
data "konnect_role" "example" {
  group_display_name = "Runtime Groups"
  display_name = "Admin"
}
```
## Argument Reference
* `group_display_name` - **(Required, String)** The display name of the role group. Must be `Runtime Groups` or `Services`
* `display_name` - **(Required, String)** The display name of the Role. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `group_name`:`name`
* `group_name` - **(String)** The name of the role group.
* `name` - **(String)** The name of the role.
* `description` - **(String)** The description of the role.
