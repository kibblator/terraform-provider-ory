package orytypes

import (
	openapiclient "github.com/ory/client-go"
)

type Hook struct {
	Config map[string]interface{} `json:"config,omitempty"`
	Hook   string                 `json:"hook,omitempty"`
}

type OryClient struct {
	APIClient *openapiclient.APIClient
	Config    *openapiclient.Project
	ProjectID string
}
