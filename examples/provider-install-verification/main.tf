terraform {
  required_providers {
    ory = {
      source = "registry.terraform.io/kibblator/ory"
    }
  }
}

provider "ory" {
}

resource "ory_registration" "example" {}
