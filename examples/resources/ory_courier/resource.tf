# Example courier settings

# Deafult
resource "ory_courier" "example" {
  delivery_strategy = "default"
}

# SMTP 
resource "ory_courier" "example" {
  delivery_strategy = "smtp"

  smtp {
    username      = "admin"
    password      = "password"
    hostname      = "smtp.example.com"
    port          = 587
    security_mode = "TLS"
    from_address  = "noreply@example.com"
    from_name     = "kibblator's Project via Ory"
    headers = {
      "X-My-Header" = "My Value"
    }
  }
}

# HTTP Basic Auth
resource "ory_courier" "example" {
  delivery_strategy = "http"

  http {
    url    = "https://ory.sh"
    method = "POST"
    body   = "SGVsbG8gV29ybGQ="
    headers = {
      "X-My-Header" = "My Value"
    }
    auth_type = "basic"
    basic_auth {
      username = "admin"
      password = "password"
    }
  }
}

# HTTP API Key
resource "ory_courier" "example" {
  delivery_strategy = "http"

  http {
    url    = "https://ory.sh"
    method = "POST"
    body   = "SGVsbG8gV29ybGQ="
    headers = {
      "X-My-Header" = "My Value"
    }
    auth_type = "api_key"
    api_key {
      name  = "my-api-key"
      value = "super-secrey-key"
      in    = "header"
    }
  }
}
