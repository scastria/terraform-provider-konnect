---
subcategory: "Identity Management"
---
# Resource: konnect_user_role
Represents a role assigned to a user to access a given entity
## Example usage
```hcl
resource "konnect_user" "User" {
  email = "Joe.Burrow@example.com"
  full_name = "Joe Burrow"
  preferred_name = "Joe"
}
data "konnect_role" "Role" {
  group_display_name = "Runtime Groups"
  display_name = "Admin"
}
resource "konnect_runtime_group" "RuntimeGroup" {
  name = "TestRG"
  description = "testing"
}
resource "konnect_user_role" "example" {
  user_id = konnect_user.User.id
  entity_id = konnect_runtime_group.RuntimeGroup.id
  entity_type_display_name = "Runtime Groups"
  entity_region = "us"
  role_display_name = data.konnect_role.Role.display_name
}
```
## Argument Reference
* `user_id` - **(Required, ForceNew, String)** The id of the user assigned the role
* `role_display_name` - **(Required, ForceNew, String)** The display name of the role.
* `entity_type_display_name` - **(Required, ForceNew, String)** The display name of the entity type, like `Runtime Groups` or `Services`.
* `entity_id` - **(Required, ForceNew, String)** The id of the entity for which the role applies.
* `entity_region` - **(Required, ForceNew, String)** The region of the entity for which the role applies.
## Attribute Reference
* `id` - **(String)** Same as `user_id`:`Guid of role assignment`
## Import
User roles can be imported using a proper value of `id` as described above
