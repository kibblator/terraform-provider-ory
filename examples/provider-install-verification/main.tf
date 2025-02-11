terraform {
  required_providers {
    ory = {
      source = "registry.terraform.io/kibblator/ory"
    }
  }
}

provider "ory" {
}

data "ory_services" "example" {}
