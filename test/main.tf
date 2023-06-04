terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

resource "konnect_authentication_settings" "AS" {
  basic_auth_enabled = true
  oidc_auth_enabled = true
  idp_mapping_enabled = false
  konnect_mapping_enabled = true
}

#resource "konnect_runtime_group" "RG" {
#  name = "asdf"
#  description = "adsfasdf"
#}
