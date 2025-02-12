package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOryRegistrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ory_registration" "test" {
  enable_registration    = false
  enable_login_hints     = true
  enable_post_signin_reg = false
  enable_password_auth   = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ory_registration.test", "enable_registration", "false"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_login_hints", "true"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_post_signin_reg", "false"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_password_auth", "true"),
					resource.TestCheckResourceAttrSet("ory_registration.test", "id"),           // Ensure the resource has an ID
					resource.TestCheckResourceAttrSet("ory_registration.test", "last_updated"), // Ensure last_updated is set
				),
			},
			// ImportState testing
			{
				ResourceName:      "ory_registration.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated", // If 'last_updated' isn't returned from the ORY API, ignore it during import verification
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "ory_registration" "test" {
  enable_registration    = true
  enable_login_hints     = false
  enable_post_signin_reg = true
  enable_password_auth   = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ory_registration.test", "enable_registration", "true"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_login_hints", "false"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_post_signin_reg", "true"),
					resource.TestCheckResourceAttr("ory_registration.test", "enable_password_auth", "false"),
					resource.TestCheckResourceAttrSet("ory_registration.test", "last_updated"), // Ensure 'last_updated' is updated
				),
			},
			// Delete testing automatically occurs in TestCase cleanup
		},
	})
}
