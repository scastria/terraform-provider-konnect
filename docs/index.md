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
  num_retries = 3
  retry_delay = 30
}
```
## Argument Reference
* `pat` - **(Required, String)** Your personal access token obtained via Konnect UI. Can be specified via env variable `KONNECT_PAT`.
* `region` - **(Optional, String)** The region for accessing region specific resources. Can be specified via env variable `KONNECT_REGION`. Allowed values: `us`, `eu`, `au`. Default: `us`
* `num_retries` - **(Optional, Integer)** Number of retries for each Konnect API call in case of 429-Too Many Requests or any 5XX status code. Can be specified via env variable `KONNECT_NUM_RETRIES`. Default: 3.
* `retry_delay` - **(Optional, Integer)** How long to wait (in seconds) in between retries. Can be specified via env variable `KONNECT_RETRY_DELAY`. Default: 30.
* `default_tags` - **(Optional, List of String)** List of tags to assign to all resources created by this provider, if the resource supports tags.
