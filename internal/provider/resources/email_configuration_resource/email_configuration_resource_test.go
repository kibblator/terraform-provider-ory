package email_configuration_resource_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/kibblator/terraform-provider-ory/internal/provider/acctest"
)

func TestAccOryEmailConfiguration_SMTP(t *testing.T) {
	randomName := acctest.GenerateRandomResourceName()
	resourceName := fmt.Sprintf("ory_email_configuration.%s", randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccOryEmailConfiguration_SMTP(randomName),
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
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated", // If 'last_updated' isn't returned from the ORY API, ignore it during import verification
				},
			},
			{
				Config: testAccOryEmailConfiguration_HTTP(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_type", "http"),
					resource.TestCheckResourceAttr(resourceName, "http_config.url", "https://ory.sh"),
					resource.TestCheckResourceAttr(resourceName, "http_config.request_method", "POST"),
					resource.TestCheckResourceAttr(resourceName, "http_config.authentication_type", "basic_auth"),
					resource.TestCheckResourceAttr(resourceName, "http_config.basic_auth.username", "username"),
					resource.TestCheckResourceAttr(resourceName, "http_config.basic_auth.password", "password"),
				),
			},
			{
				Config: testAccOryEmailConfiguration_Default(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_type", "default"),
				),
			},
		},
	})
}

func testAccOryEmailConfiguration_SMTP(randomName string) string {
	return fmt.Sprintf(`
resource "ory_email_configuration" "%s" {
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
`, randomName)
}

func testAccOryEmailConfiguration_Default(randomName string) string {
	return fmt.Sprintf(`
resource "ory_email_configuration" "%s" {
  server_type = "default"
}
`, randomName)
}

func testAccOryEmailConfiguration_HTTP(randomName string) string {
	return fmt.Sprintf(`
resource "ory_email_configuration" "%s" {
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
`, randomName)
}
