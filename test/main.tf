terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

data "konnect_user" "U" {
  full_name = "Travis Valentine"
}

#resource "konnect_identity_provider" "IP" {
#  issuer = "https://greenst.okta.com/oauth2/default"
#  client_id = "0oambh387v9ETDgCz2p7"
#  login_path = "gsa"
#  scopes = [
#    "email",
#    "openid",
#    "profile"
#  ]
#  email_claim_mapping = "email"
#  name_claim_mapping = "name"
#  groups_claim_mapping = "groups"
#}

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
