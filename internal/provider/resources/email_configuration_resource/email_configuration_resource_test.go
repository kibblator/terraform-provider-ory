package email_configuration_resource_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/kibblator/terraform-provider-ory/internal/provider/acctest"
)

func TestAccOryEmailConfiguration_Default(t *testing.T) {
	resourceName := "ory_email_configuration.default"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOryEmailConfiguration_Default(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_type", "default"),
				),
			},
		},
	})
}

func testAccOryEmailConfiguration_Default() string {
	return `
resource "ory_email_configuration" "default" {
  server_type = "default"
}
`
}

func TestAccOryEmailConfiguration_SMTP(t *testing.T) {
	resourceName := "ory_email_configuration.smtp"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOryEmailConfiguration_SMTP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_type", "smtp"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.sender_name", "Ory"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.sender_address", "noreply@examplecompany.com"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.host", "smtp.examplecompany.com"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.port", "587"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.security", "starttls"),
					resource.TestCheckResourceAttr(resourceName, "smtp_config.username", "username"),
				),
			},
		},
	})
}

func testAccOryEmailConfiguration_SMTP() string {
	return `
resource "ory_email_configuration" "smtp" {
  server_type = "smtp"

  smtp_config = {
    sender_name    = "Ory"
    sender_address = "noreply@examplecompany.com"
    host           = "smtp.examplecompany.com"
    port           = "587"
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
`
}

func TestAccOryEmailConfiguration_HTTP(t *testing.T) {
	resourceName := "ory_email_configuration.http"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOryEmailConfiguration_HTTP(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_type", "http"),
					resource.TestCheckResourceAttr(resourceName, "http_config.url", "https://ory.sh"),
					resource.TestCheckResourceAttr(resourceName, "http_config.request_method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "http_config.authentication_type", "basic_auth"),
					resource.TestCheckResourceAttr(resourceName, "http_config.basic_auth.username", "username"),
					resource.TestCheckResourceAttr(resourceName, "http_config.basic_auth.password", "password"),
				),
			},
		},
	})
}

func testAccOryEmailConfiguration_HTTP() string {
	return `
resource "ory_email_configuration" "http" {
  server_type = "http"

  http_config = {
    url                 = "https://ory.sh"
    request_method      = "POST"
    authentication_type = "basic_auth"

    basic_auth = {
      username = "username"
      password = "password"
    }

    action_body = "aGVsbG8gd29ybGQ="
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

  lifecycle {
    ignore_changes = [
      http_config.action_body
    ]
  }
}
`
}
