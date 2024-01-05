---
subcategory: "Identity Management"
---
# Data Source: konnect_role
Represents a role
## Example usage
```hcl
data "konnect_role" "example" {
  entity_type_display_name = "Control Planes"
  display_name = "Admin"
}
```
## Argument Reference
* `entity_type_display_name` - **(Required, String)** The display name of the role entity type. Must be `Control Planes` or `Services`
* `display_name` - **(Required, String)** The display name of the Role. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `entity_type_name`:`name`
* `entity_type_name` - **(String)** The name of the role entity type.
* `name` - **(String)** The name of the role.
* `description` - **(String)** The description of the role.
