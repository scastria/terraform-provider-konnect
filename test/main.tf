terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

data "konnect_team_role" "TR" {
  team_id = data.konnect_team.T.id
  entity_type_display_name = "Runtime Groups"
}

#data "konnect_role" "R" {
#  group_display_name = "Runtime Groups"
#  display_name = "Admin"
#}

#resource "konnect_team" "T" {
#  name = "ShawnTest"
#  description = "testing"
#}

data "konnect_team" "T" {
  name = "runtime-admin"
}

#resource "konnect_user" "U" {
#  email = "jblow@example.com"
#  full_name = "Joe Blow"
#  preferred_name = "Joe"
#}

#data "konnect_user" "U" {
#  search_full_name = "Julia"
#}

#resource "konnect_team_user" "TU" {
#  team_id = konnect_team.T.id
#  user_id = data.konnect_user.U.id
#}

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
