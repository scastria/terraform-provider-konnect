---
subcategory: "Identity Management"
---
# Resource: konnect_team_mappings
Represents the mappings between an external identity provider group and a Konnect team
## Example usage
```hcl
resource "konnect_team" "Team" {
  name = "Test"
  description = "testing"
}
resource "konnect_team_mappings" "example" {
  mapping {
    group = "external IdP group"
    team_ids = [
      data.konnect_team.Team.id
    ]
  }
}
```
## Argument Reference
* `mapping` - **(Optional, set{mapping})** Configuration block for a mapping.  Can be specified multiple times for each mapping.  Each block supports the fields documented below.
### mapping
* `group` - **(Required, String)** Identifier of an IdP group that is contained with OIDC ID token for groups claim
* `team_ids` - **(Required, List of String)** Konnect teams that should map to this group.
## Attribute Reference
* `id` - **(String)** Always equal to `team-mappings`
## Import
Team mappings can be imported using a proper value of `id` as described above
