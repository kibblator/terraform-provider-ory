terraform {
  required_providers {
    ory = {
      source = "registry.terraform.io/kibblator/ory"
    }
  }
}


provider "ory" {
}

resource "ory_registration" "reg_settings" {
  enable_registration    = false
  enable_login_hints     = false
  enable_post_signin_reg = false
  enable_password_auth   = false
}


output "registration_settings" {
  value = ory_registration.reg_settings
}
