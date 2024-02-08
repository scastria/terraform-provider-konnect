# Kong Konnect Provider
The Kong Konnect provider is used to interact with the Kong Konnect API.  The provider
needs to be configured with the proper credentials before it can be used.

This provider does NOT cover 100% of the Konnect API.  If there is something missing
that you would like to be added, please submit an Issue in corresponding GitHub repo.
## Example Usage
```hcl
terraform {
  required_providers {
    konnect = {
      source  = "scastria/konnect"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Konnect Provider
provider "konnect" {
  pat = "XXXX"
}
```
## Argument Reference
* `pat` - **(Required, String)** Your personal access token obtained via Konnect UI. Can be specified via env variable `KONNECT_PAT`.
* `region` - **(Optional, String)** The region for accessing region specific resources. Can be specified via env variable `KONNECT_REGION`. Allowed values: `us`, `eu`, `au`. Default: `us`
