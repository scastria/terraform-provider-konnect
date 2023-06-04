terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

resource "konnect_identity_provider" "IP" {
  issuer = "https://greenst.okta.com/oauth2/default"
  client_id = "0oambh387v9ETDgCz2p7"
  login_path = "gsa"
  scopes = [
    "email",
    "openid",
    "profile"
  ]
  claim_mappings = {
    email = "email"
    groups = "groups"
    name = "name"
  }
}

#resource "konnect_authentication_settings" "AS" {
#  basic_auth_enabled = true
#  oidc_auth_enabled = true
#  idp_mapping_enabled = false
#  konnect_mapping_enabled = true
#}

#resource "konnect_runtime_group" "RG" {
#  name = "asdf"
#  description = "adsfasdf"
#}
