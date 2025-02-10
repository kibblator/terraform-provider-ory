terraform {
  required_providers {
    ory = {
      source = "registry.terraform.io/kibblator/ory"
    }
  }
}

provider "ory" {
  workspace_api_key = "ory_wak_1234567890"
  project_id        = "project-guid-here"
}

data "ory_services" "example" {}
