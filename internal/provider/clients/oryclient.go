package oryclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
	"github.com/ory/client-go"
)

type Client struct {
	BaseURL    string
	APIKey     string
	ProjectID  string
	HTTPClient *http.Client
}

type OryClient struct {
	APIClient     *Client
	ProjectConfig *orytypes.Project
	ProjectID     string
	Mutex         sync.Mutex
}

func NewClient(baseUrl, apiKey, projectID string) *Client {
	return &Client{
		BaseURL:    fmt.Sprintf("https://%s", baseUrl),
		APIKey:     apiKey,
		ProjectID:  projectID,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetProject(m *sync.Mutex) (*orytypes.Project, error) {
	m.Lock()
	defer m.Unlock()

	url := fmt.Sprintf("%s/projects/%s", c.BaseURL, c.ProjectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch project: %s", body)
	}

	var config orytypes.Project
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Client) PatchProject(revisionID string, patchData []client.JsonPatch, m *sync.Mutex) (*orytypes.ProjectConfig, error) {
	m.Lock()
	defer m.Unlock()

	url := fmt.Sprintf("%s/projects/%s/revision/%s", c.BaseURL, c.ProjectID, revisionID)
	data, err := json.Marshal(patchData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to patch project: %s", body)
	}

	var updatedConfig orytypes.ProjectConfig
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &updatedConfig)
	if err != nil {
		return nil, err
	}

	return &updatedConfig, nil
}
