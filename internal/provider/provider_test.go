package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	providerConfig = `
provider "ory" {
  workspace_api_key    = "` + os.Getenv("ORY_WORKSPACE_API_KEY") + `"
  project_id = "` + os.Getenv("ORY_PROJECT_ID") + `"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ory": providerserver.NewProtocol6WithError(New("test")()),
	}
)
