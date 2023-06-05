---
subcategory: "Identity Management"
---
# Resource: konnect_team_user
Represents a member of a team
## Example usage
```hcl
resource "konnect_team" "Team" {
  name = "Test"
  description = "testing"
}
data "konnect_user" "User" {
  search_full_name = "Joe"
}
resource "konnect_team_user" "example" {
  team_id = konnect_team.Team.id
  user_id = data.konnect_user.User.id
}
```
## Argument Reference
* `team_id` - **(Required, ForceNew, String)** The id of the team.
* `user_id` - **(Required, ForceNew, String)** The id of the user.
## Attribute Reference
* `id` - **(String)** Same as `team_id`:`user_id`
## Import
Team users can be imported using a proper value of `id` as described above
