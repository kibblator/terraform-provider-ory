package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type servicesDataSource struct {
}

// NewServicesDataSource is a helper function to simplify the provider implementation.
func NewServicesDataSource() datasource.DataSource {
	return &servicesDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

// Metadata returns the data source type name.
func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

// Schema defines the schema for the data source.
func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Read refreshes the Terraform state with the latest data.
func (d *servicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
