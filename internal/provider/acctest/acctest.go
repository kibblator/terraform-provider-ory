package acctest

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/kibblator/terraform-provider-ory/internal/provider"
)

var ProviderConfig = `
provider "ory" {
  workspace_api_key    = "` + os.Getenv("ORY_WORKSPACE_API_KEY") + `"
  project_id = "` + os.Getenv("ORY_PROJECT_ID") + `"
}
`

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ory": providerserver.NewProtocol6WithError(provider.New("test")()),
}
