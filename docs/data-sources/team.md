# Data Source: konnect_team
Represents a Konnect team
## Example usage
```hcl
data "konnect_team" "example" {
  search_name = "Panther"
}
```
## Argument Reference
* `search_name` - **(Optional, String)** The search string to apply to the name of the team. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the team. Uses equality.
## Attribute Reference
* `id` - **(String)** Guid
* `description` - **(String)** The preferred description of the team.
* `is_predefined` - **(Boolean)** Whether the team is predefined.
