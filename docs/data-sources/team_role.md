# Data Source: konnect_team_role
Represents a role assigned to a team
## Example usage
```hcl
data "konnect_team" "Team" {
  name = "runtime-admin"
}
data "konnect_team_role" "example" {
  team_id = data.konnect_team.Team.id
  entity_type_display_name = "Runtime Groups"
}
```
## Argument Reference
* `team_id` - **(Required, String)** The id of the team assigned the role
* `search_role_display_name` - **(Optional, String)** The search string to apply to the display name of the role. Uses contains.
* `role_display_name` - **(Optional, String)** The filter string to apply to the display name of the role. Uses equality.
* `search_entity_type_display_name` - **(Optional, String)** The search string to apply to the display name of the entity type, like `Runtime Groups` or `Services`. Uses contains.
* `entity_type_display_name` - **(Optional, String)** The filter string to apply to the display name of the entity type, like `Runtime Groups` or `Services`. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `team_id`:`Guid of role assignment`
* `entity_id` - **(String)** The id of the entity for which the role applies.
* `entity_region` - **(String)** The region of the entity for which the role applies.
