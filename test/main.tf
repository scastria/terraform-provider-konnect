terraform {
  required_providers {
    konnect = {
      source = "github.com/scastria/konnect"
    }
  }
}

provider "konnect" {
}

#data "msgraph_group" "MyGroup" {
#  display_name = "ShawnTest"
##  wait_until_exists = true
##  wait_timeout = 55
##  wait_polling_interval = 3
#}
#
#output "MyOutput" {
#  value = data.msgraph_group.MyGroup.id
#}
