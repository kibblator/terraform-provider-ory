package acctest

import (
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/kibblator/terraform-provider-ory/internal/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"ory": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6(provider.New("dev")())(), nil
	},
}

func GenerateRandomResourceName() string {
	const resourceNameLength = 10
	const charset = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder
	for i := 0; i < resourceNameLength; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func TestAccPreCheck(t *testing.T) {
	TestAccPreCheck_Provider(t)
}

func TestAccPreCheck_Provider(t *testing.T) {
	host := os.Getenv("ORY_HOST")
	project_id := os.Getenv("ORY_PROJECT_ID")
	workspace_api_key := os.Getenv("ORY_WORKSPACE_API_KEY")

	if host == "" && project_id == "" && workspace_api_key == "" {
		t.Fatal("Provider environment variables need to be setup for this test")
	}
}
