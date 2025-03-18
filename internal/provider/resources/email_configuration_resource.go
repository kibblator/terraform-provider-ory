package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
)

type emailConfigurationResource struct {
	oryClient *orytypes.OryClient
}

var (
	_ resource.Resource                = &emailConfigurationResource{}
	_ resource.ResourceWithConfigure   = &emailConfigurationResource{}
	_ resource.ResourceWithImportState = &emailConfigurationResource{}
)

// Configure implements resource.ResourceWithConfigure.
func (r *emailConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*orytypes.OryClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected OryClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.oryClient = client
}

// Metadata returns the resource type name.
func (r *emailConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_configuration"
}

// Schema implements resource.Resource.
func (r *emailConfigurationResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"server_type": schema.StringAttribute{
				Description: "The type of the email server.",
				Required:    true,
			},
			"smtp_config": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sender_name": schema.StringAttribute{
							Description: "The name of the sender.",
							Required:    true,
						},
						"sender_address": schema.StringAttribute{
							Description: "The email address of the sender.",
							Required:    true,
						},
						"host": schema.StringAttribute{
							Description: "The SMTP server host.",
							Required:    true,
						},
						"port": schema.NumberAttribute{
							Description: "The SMTP server port.",
							Required:    true,
						},
					},
			},
		},
	},
}

// Create implements resource.Resource.
func (r *emailConfigurationResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Read implements resource.Resource.
func (r *emailConfigurationResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (r *emailConfigurationResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (r *emailConfigurationResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

// ImportState implements resource.ResourceWithImportState.
func (r *emailConfigurationResource) ImportState(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	panic("unimplemented")
}
