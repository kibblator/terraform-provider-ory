# Example registration settings
resource "ory_registration" "example" {
  enable_registration    = true
  enable_login_hints     = false
  enable_post_signin_reg = true
  enable_password_auth   = false
}
