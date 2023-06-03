terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

data "konnect_runtime_group" "RG" {
  name = "development"
}

output "MyOutput" {
  value = data.konnect_runtime_group.RG
}
