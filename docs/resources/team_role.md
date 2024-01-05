---
subcategory: "Identity Management"
---
# Resource: konnect_team_role
Represents a role assigned to a team to access a given entity
## Example usage
```hcl
resource "konnect_team" "Team" {
  name = "Test"
  description = "testing"
}
data "konnect_role" "Role" {
  group_display_name = "Control Planes"
  display_name = "Admin"
}
resource "konnect_control_plane" "ControlPlane" {
  name = "TestRG"
  description = "testing"
}
resource "konnect_team_role" "example" {
  team_id = konnect_team.Team.id
  entity_id = konnect_control_plane.ControlPlane.id
  entity_type_display_name = "Control Planes"
  entity_region = "us"
  role_display_name = data.konnect_role.Role.display_name
}
```
## Argument Reference
* `team_id` - **(Required, ForceNew, String)** The id of the team assigned the role
* `role_display_name` - **(Required, ForceNew, String)** The display name of the role.
* `entity_type_display_name` - **(Required, ForceNew, String)** The display name of the entity type, like `Control Planes` or `Services`.
* `entity_id` - **(Required, ForceNew, String)** The id of the entity for which the role applies.
* `entity_region` - **(Required, ForceNew, String)** The region of the entity for which the role applies.
## Attribute Reference
* `id` - **(String)** Same as `team_id`:`Guid of role assignment`
## Import
Team roles can be imported using a proper value of `id` as described above
