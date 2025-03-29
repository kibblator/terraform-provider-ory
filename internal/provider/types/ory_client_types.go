package orytypes

import (
	openapiclient "github.com/ory/client-go"
)

type OryClient struct {
	APIClient     *openapiclient.APIClient
	ProjectConfig *openapiclient.Project
	ProjectID     string
}
