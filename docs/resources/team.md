---
subcategory: "Identity Management"
---
# Resource: konnect_team
Represents a team
## Example usage
```hcl
resource "konnect_team" "example" {
  name = "Panthers"
  description = "dev team"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the team.
* `description` - **(Optional, String)** The description of the team.
## Attribute Reference
* `id` - **(String)** Guid
* `is_predefined` - **(Boolean)** Whether the team is predefined.
## Import
Teams can be imported using a proper value of `id` as described above
