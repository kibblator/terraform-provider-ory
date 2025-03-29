package registration_resource_test

import (
	"fmt"
	"testing"

	"github.com/kibblator/terraform-provider-ory/internal/provider/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOryRegistrationResource(t *testing.T) {
	t.Parallel()

	randomName := acctest.GenerateRandomResourceName()
	resourceName := fmt.Sprintf("ory_registration.%s", randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: fmt.Sprintf(`
resource "ory_registration" "%s" {
  enable_registration    = false
  enable_login_hints     = true
  enable_post_signin_reg = false
  enable_password_auth   = true
}
`, randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_registration", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_login_hints", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_post_signin_reg", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_password_auth", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),           // Ensure the resource has an ID
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"), // Ensure last_updated is set
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_updated", // If 'last_updated' isn't returned from the ORY API, ignore it during import verification
				},
			},
			// Update and Read testing
			{
				Config: fmt.Sprintf(`
			resource "ory_registration" "%s" {
			  enable_registration    = true
			  enable_login_hints     = false
			  enable_post_signin_reg = true
			  enable_password_auth   = false
			}
			`, randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_registration", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_login_hints", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_post_signin_reg", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_password_auth", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"), // Ensure 'last_updated' is updated
				),
			},
			// Delete testing automatically occurs in TestCase cleanup
		},
	})
}
