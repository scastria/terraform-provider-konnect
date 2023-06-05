---
subcategory: "Identity Management"
---
# Data Source: konnect_user
Represents a user
## Example usage
```hcl
data "konnect_user" "example" {
  search_email = "@example.com"
}
```
## Argument Reference
* `search_email` - **(Optional, String)** The search string to apply to the email of the user. Uses contains.
* `email` - **(Optional, String)** The filter string to apply to the email of the user. Uses equality.
* `search_full_name` - **(Optional, String)** The search string to apply to the full name of the user. Uses contains.
* `full_name` - **(Optional, String)** The filter string to apply to the full name of the user. Uses equality.
* `active` - **(Optional, Boolean)** The filter flag to apply to the active flag of the user. Uses equality. Default: `true`
## Attribute Reference
* `id` - **(String)** Guid
* `preferred_name` - **(String)** The preferred name of the user.
