# Email configuration using default server type
resource "ory_email_configuration" "example" {
  server_type = "default"
}

# Email configuration using smtp server type
resource "ory_email_configuration" "example" {
  server_type = "smtp"

  smtp_config = {
    sender_name    = "Ory"
    sender_address = "noreply@examplecompany.com"
    host           = "smtp.examplecompany.com"
    port           = 587
    security       = "starttls"
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
resource "ory_email_configuration" "example" {
  server_type = "http"

  http_config = {
    url                 = "https://ory.sh"
    request_method      = "POST"                                     # Can be GET, POST, PUT, PATCH
    authentication_type = "none"                                     # Can be none, basic, apikey
    username            = "username"                                 # ignored if authentication is not basic
    password            = "password"                                 # ignored if authentication is not basic
    api_key             = "api_key"                                  # ignored if authentication is not apikey
    transport_mode      = "header"                                   # Can be header, cookie - ignored if authentication is not apikey
    action_body         = base64encode(file("path/to/body.jsonnet")) # base64 encoded string containing the jsonnet body (uses default payload from ory docs if not provided)
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
