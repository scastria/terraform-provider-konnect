terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

#resource "konnect_runtime_group" "RG" {
#  name = "asdf"
#  description = "adsfasdf"
#}
