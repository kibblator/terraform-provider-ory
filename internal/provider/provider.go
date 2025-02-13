package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openapiclient "github.com/ory/client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &oryProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &oryProvider{
			version: version,
		}
	}
}

// oryProvider is the provider implementation.
type oryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// oryProvider maps provider schema data to a Go type.
type oryProviderModel struct {
	Host            types.String `tfsdk:"host"`
	ProjectId       types.String `tfsdk:"project_id"`
	WorkSpaceApiKey types.String `tfsdk:"workspace_api_key"`
}

type OryClient struct {
	APIClient *openapiclient.APIClient
	Config    *openapiclient.Project
	ProjectID string
}

// Metadata returns the provider type name.
func (p *oryProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ory"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *oryProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for the Ory Network console API. May also be provided with the ORY_HOST environment variable.",
				Optional:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "The project ID for the target Ory Network Project. May also be provided with the ORY_PROJECT_ID environment variable.",
				Optional:    true,
			},
			"workspace_api_key": schema.StringAttribute{
				Description: "Your Ory Network workspace API key. May also be provided with the ORY_WORKSPACE_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares an Ory API client for data sources and resources.
func (p *oryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Ory client")
	// Retrieve provider data from configuration
	var config oryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Ory API Host",
			"The provider cannot create the Ory API client as there is an unknown configuration value for the Ory API host.",
		)
	}

	if config.ProjectId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Unknown Ory API Project Id",
			"The provider cannot create the Ory API client as there is an unknown configuration value for the Ory API project id.",
		)
	}

	if config.WorkSpaceApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_api_key"),
			"Unknown Ory API Workspace API Key",
			"The provider cannot create the Ory API client as there is an unknown configuration value for the Ory API workspace API key.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("ORY_HOST")
	project_id := os.Getenv("ORY_PROJECT_ID")
	workspace_api_key := os.Getenv("ORY_WORKSPACE_API_KEY")

	tflog.Debug(ctx, "Checking environment variables for Ory configuration", map[string]interface{}{
		"ory_host":              host,
		"ory_project_id":        project_id,
		"ory_workspace_api_key": workspace_api_key,
	})

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ProjectId.IsNull() {
		project_id = config.ProjectId.ValueString()
	}

	if !config.WorkSpaceApiKey.IsNull() {
		workspace_api_key = config.WorkSpaceApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		host = "api.console.ory.sh"
	}

	if project_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Missing Ory API Project Id",
			"The provider cannot create the Ory API client as there is a missing or empty value for the Ory API project id. "+
				"If this is already set, ensure the value is not empty.",
		)
	}

	if workspace_api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("workspace_api_key"),
			"Missing Ory API Workspace API Key",
			"The provider cannot create the Ory API client as there is a missing or empty value for the Ory API workspace API key. "+
				"If this is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ory_host", host)
	ctx = tflog.SetField(ctx, "ory_project_id", project_id)
	ctx = tflog.SetField(ctx, "ory_workspace_api_key", workspace_api_key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ory_workspace_api_key")

	tflog.Debug(ctx, "Creating Ory client")

	// Create a new Ory client using the configuration values and pull configuration
	configuration := openapiclient.NewConfiguration()
	configuration.Host = host
	configuration.AddDefaultHeader("Authorization", "Bearer "+workspace_api_key)

	apiClient := openapiclient.NewAPIClient(configuration)
	response, r, err := apiClient.ProjectAPI.GetProject(ctx, project_id).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ProjectAPI.GetProject``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		resp.Diagnostics.AddError(
			"Unable to get project configuration using the Ory API",
			"An unexpected error occurred when calling GetProject on the Ory API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Ory Client Error: "+err.Error(),
		)
		return
	}

	client := &OryClient{
		APIClient: apiClient,
		Config:    response,
		ProjectID: project_id,
	}

	// Make the Ory config available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Ory client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *oryProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *oryProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRegistrationResource,
	}
}
