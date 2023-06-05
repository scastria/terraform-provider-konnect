# Resource: konnect_user
Represents a user
## Example usage
```hcl
resource "konnect_user" "example" {
  email = "Joe.Burrow@example.com"
  full_name = "Joe Burrow"
  preferred_name = "Joe"
}
```
## Argument Reference
* `email` - **(Required, String)** The email of the user.
* `full_name` - **(Optional, String)** The full name of the user.
* `preferred_name` - **(Optional, String)** The preferred name of the user.
## Attribute Reference
* `id` - **(String)** Guid
* `active` - **(Boolean)** Whether the user is active.
## Import
Users can be imported using a proper value of `id` as described above
