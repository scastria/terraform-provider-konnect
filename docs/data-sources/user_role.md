---
subcategory: "Identity Management"
---
# Data Source: konnect_user_role
Represents a role assigned to a user
## Example usage
```hcl
data "konnect_user" "User" {
  search_full_name = "Joe"
}
data "konnect_user_role" "example" {
  user_id = data.konnect_user.User.id
  entity_type_display_name = "Runtime Groups"
}
```
## Argument Reference
* `user_id` - **(Required, String)** The id of the user assigned the role
* `search_role_display_name` - **(Optional, String)** The search string to apply to the display name of the role. Uses contains.
* `role_display_name` - **(Optional, String)** The filter string to apply to the display name of the role. Uses equality.
* `search_entity_type_display_name` - **(Optional, String)** The search string to apply to the display name of the entity type, like `Runtime Groups` or `Services`. Uses contains.
* `entity_type_display_name` - **(Optional, String)** The filter string to apply to the display name of the entity type, like `Runtime Groups` or `Services`. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `user_id`:`Guid of role assignment`
* `entity_id` - **(String)** The id of the entity for which the role applies.
* `entity_region` - **(String)** The region of the entity for which the role applies.
