#Email configuration using default server type
resource "ory_email_configuration" "default" {
  server_type = "default"
}

#Email configuration using smtp server type
resource "ory_email_configuration" "smtp" {
  server_type = "smtp"

  smtp_config = {
    sender_name    = "Ory"
    sender_address = "noreply@examplecompany.com"
    host           = "smtp.examplecompany.com"
    port           = "587"
    security       = "starttls" # options are starttls, starttls_notrust, cleartext, implicittls, implicittls_notrust
    username       = "username"
    password       = "password"
  }

  smtp_headers = [
    {
      key   = "X-Header-1"
      value = "value-1"
    },
    {
      key   = "X-Header-2"
      value = "value-2"
    }
  ]
}

# Email configuration using http server type
resource "ory_email_configuration" "http" {
  server_type = "http"

  http_config = {
    url                 = "https://ory.sh"
    request_method      = "POST"       # Can be GET, POST, PUT, PATCH
    authentication_type = "basic_auth" # Can be none, basic_auth, api_key

    basic_auth = { # set only if authentication_type is basic_auth
      username = "username"
      password = "password"
    }

    api_key = {                 # # set only if authentication_type is api_key
      transport_mode = "header" # Can be header, cookie
      name           = "x-api-key"
      value          = "super secret value"
    }

    action_body = base64encode(file("./somefile.jsonnet")) # base64 encoded string containing the jsonnet body (uses default payload from ory docs if not provided)
  }

  smtp_headers = [
    {
      key   = "X-Header-1"
      value = "value-1"
    },
    {
      key   = "X-Header-2"
      value = "value-2"
    }
  ]
}
