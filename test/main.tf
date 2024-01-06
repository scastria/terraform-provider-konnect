terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

#resource "konnect_plugin" "P" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  name = "rate-limiting"
#  protocols = [
#    "grpc",
#    "grpcs",
#    "http",
#    "https"
#  ]
#  config_json = <<EOF
#{
#  "minute": 8,
#  "second": 7
#}
#EOF
#}

#resource "konnect_consumer" "C" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  username = "Shawn"
#  custom_id = "Bob"
#}

#resource "konnect_service" "S" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  host = "mockbin.org"
#  name = "TFTest"
#}

#resource "konnect_route" "R" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  name = "TFRoute"
#  protocols = ["http"]
#  methods = ["GET"]
#  paths = ["/tf"]
#  service_id = konnect_service.S.service_id
#  header {
#    name = "sear"
#    values = ["kevin"]
#  }
#}

#data "konnect_nodes" "Ns" {
#  control_plane_id = data.konnect_control_plane.RG.id
#}

#resource "konnect_user_role" "UR" {
#  user_id = data.konnect_user.U.id
#  entity_id = data.konnect_control_plane.RG.id
#  entity_type_display_name = "Control Planes"
#  entity_region = "us"
#  role_display_name = data.konnect_role.R.display_name
#}

#data "konnect_team_role" "TR" {
#  team_id = data.konnect_team.T.id
#  entity_type_display_name = "Control Planes"
#}

#data "konnect_role" "R" {
#  group_display_name = "Control Planes"
#  display_name = "Admin"
#}

#resource "konnect_team" "T" {
#  name = "ShawnTest"
#  description = "testing"
#}

#resource "konnect_team_role" "TR" {
#  team_id = konnect_team.T.id
#  entity_id = konnect_control_plane.RG.id
#  entity_type_display_name = "Control Planes"
#  entity_region = "us"
#  role_display_name = data.konnect_role.R.display_name
#}

#data "konnect_team" "T" {
#  name = "runtime-admin"
#}

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

#resource "konnect_control_plane" "RG" {
#  name = "ShawnRG"
#  description = "testing"
#}

data "konnect_control_plane" "RG" {
  name = "development"
}

data "konnect_consumer" "C" {
  control_plane_id = data.konnect_control_plane.RG.id
  search_username = "web"
#  search_custom_id = null
}

#resource "konnect_consumer_key" "CK" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  consumer_id = data.konnect_consumer.C.consumer_id
#  key = "shawn"
#}
#resource "konnect_consumer_acl" "CACL" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  consumer_id = data.konnect_consumer.C.consumer_id
#  group = "shawn2"
#}
#resource "konnect_consumer_basic" "CB" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  consumer_id = data.konnect_consumer.C.consumer_id
#  username = "shawn"
#  password = "pass2"
#}
#resource "konnect_consumer_hmac" "CHMAC" {
#  control_plane_id = data.konnect_control_plane.RG.id
#  consumer_id      = data.konnect_consumer.C.consumer_id
#  username         = "shawn"
##  secret           = "pass2"
#}
