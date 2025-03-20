package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
)

var (
	_ resource.Resource                   = &emailConfigurationResource{}
	_ resource.ResourceWithConfigure      = &emailConfigurationResource{}
	_ resource.ResourceWithImportState    = &emailConfigurationResource{}
	_ resource.ResourceWithValidateConfig = &emailConfigurationResource{}
)

type emailConfigurationResource struct {
	oryClient *orytypes.OryClient
}

type SMTPConfig struct {
	SenderName    *string `tfsdk:"sender_name"`
	SenderAddress *string `tfsdk:"sender_address"`
	Host          *string `tfsdk:"host"`
	Port          *int    `tfsdk:"port"`
	Security      *string `tfsdk:"security"`
	Username      *string `tfsdk:"username"`
	Password      *string `tfsdk:"password"`
}

type HTTPConfig struct {
	Url                *string `tfsdk:"url"`
	RequestMethod      *string `tfsdk:"request_method"`
	AuthenticationType *string `tfsdk:"authentication_type"`
	Username           *string `tfsdk:"username"`
	Password           *string `tfsdk:"password"`
	ApiKey             *string `tfsdk:"api_key"`
	TransportMode      *string `tfsdk:"transport_mode"`
	ActionBody         *string `tfsdk:"action_body"`
}

type SMTPHeader struct {
	Key   string `tfsdk:"key"`
	Value string `tfsdk:"value"`
}

type emailConfigurationResourceModel struct {
	ServerType  string       `tfsdk:"server_type"`
	SMTPConfig  *SMTPConfig  `tfsdk:"smtp_config"`
	HTTPConfig  *HTTPConfig  `tfsdk:"http_config"`
	SMTPHeaders []SMTPHeader `tfsdk:"smtp_headers"`
}

func NewEmailConfigurationResource() resource.Resource {
	return &emailConfigurationResource{}
}

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
			"smtp_config": schema.SingleNestedAttribute{
				Description: "SMTP configuration block (optional, but fields required if present).",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sender_name": schema.StringAttribute{
						Description: "The name of the sender.",
						Optional:    true,
					},
					"sender_address": schema.StringAttribute{
						Description: "The email address of the sender.",
						Optional:    true,
					},
					"host": schema.StringAttribute{
						Description: "The SMTP server host.",
						Optional:    true,
					},
					"port": schema.NumberAttribute{
						Description: "The SMTP server port.",
						Optional:    true,
					},
					"security": schema.StringAttribute{
						Description: "The security type of the SMTP server.",
						Optional:    true,
					},
					"username": schema.StringAttribute{
						Description: "The username for the SMTP server.",
						Optional:    true,
					},
					"password": schema.StringAttribute{
						Description: "The password for the SMTP server.",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
			"http_config": schema.SingleNestedAttribute{
				Description: "HTTP configuration block (optional, but fields required if present).",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "The URL of the HTTP server.",
						Optional:    true,
					},
					"request_method": schema.StringAttribute{
						Description: "The request method for the HTTP server.",
						Optional:    true,
					},
					"authentication_type": schema.StringAttribute{
						Description: "The authentication type for the HTTP server.",
						Optional:    true,
					},
					"username": schema.StringAttribute{
						Description: "The username for the HTTP server.",
						Optional:    true,
					},
					"password": schema.StringAttribute{
						Description: "The password for the HTTP server.",
						Optional:    true,
						Sensitive:   true,
					},
					"api_key": schema.StringAttribute{
						Description: "The API key for the HTTP server.",
						Optional:    true,
						Sensitive:   true,
					},
					"transport_mode": schema.StringAttribute{
						Description: "The transport mode for the HTTP server.",
						Optional:    true,
					},
					"action_body": schema.StringAttribute{
						Description: "The action body for the HTTP server.",
						Optional:    true,
					},
				},
			},
			"smtp_headers": schema.ListNestedAttribute{
				Description: "SMTP headers block (required when server_type is smtp or http).",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Description: "The key of the SMTP header.",
							Required:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value of the SMTP header.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func (r emailConfigurationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data emailConfigurationResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServerType == "default" {
		if data.SMTPConfig != nil || data.HTTPConfig != nil {
			resp.Diagnostics.AddError("smtp_config and http_config are not allowed with server_type default", "SMTP and HTTP configurations are not allowed with default server type.")
		}
	}

	if data.ServerType == "smtp" {
		if data.HTTPConfig != nil {
			resp.Diagnostics.AddError("http_config is not allowed with server_type smtp", "HTTP configuration is not allowed with SMTP server type.")
			return
		}

		if data.SMTPConfig == nil {
			resp.Diagnostics.AddError("smtp_config is missing", "SMTP configuration is required with SMTP server type.")
			return
		}

		if data.SMTPConfig.SenderName == nil || *data.SMTPConfig.SenderName == "" {
			resp.Diagnostics.AddError("sender_name is missing from smtp_config", "Sender name is required for SMTP server type.")
		}

		if data.SMTPConfig.SenderAddress == nil || *data.SMTPConfig.SenderAddress == "" {
			resp.Diagnostics.AddError("sender_address is missing from smtp_config", "Sender address is required for SMTP server type.")
		}

		if data.SMTPConfig.Host == nil || *data.SMTPConfig.Host == "" {
			resp.Diagnostics.AddError("host is missing from smtp_config", "Host is required for SMTP server type.")
		}

		if data.SMTPConfig.Port == nil || *data.SMTPConfig.Port == 0 {
			resp.Diagnostics.AddError("port is missing from smtp_config", "Port is required for SMTP server type.")
		}

		if data.SMTPConfig.Security == nil || *data.SMTPConfig.Security == "" {
			resp.Diagnostics.AddError("security is missing from smtp_config", "Security is required for SMTP server type.")
		}

		if data.SMTPConfig.Username == nil || *data.SMTPConfig.Username == "" {
			resp.Diagnostics.AddError("username is missing from smtp_config", "Username is required for SMTP server type.")
		}

		if data.SMTPConfig.Password == nil || *data.SMTPConfig.Password == "" {
			resp.Diagnostics.AddError("password is missing from smtp_config", "Password is required for SMTP server type.")
		}
	}

	if data.ServerType == "http" {
		if data.SMTPConfig != nil {
			resp.Diagnostics.AddError("smtp_config is not allowed with server_type http", "SMTP configuration is not allowed with HTTP server type.")
			return
		}

		if data.HTTPConfig == nil {
			resp.Diagnostics.AddError("http_config is missing", "HTTP configuration is required with HTTP server type.")
			return
		}

		if data.HTTPConfig.Url == nil || *data.HTTPConfig.Url == "" {
			resp.Diagnostics.AddError("url is missing from http_config", "URL is required for HTTP server type.")
		}

		if data.HTTPConfig.RequestMethod == nil || *data.HTTPConfig.RequestMethod == "" {
			resp.Diagnostics.AddError("request_method is missing from http_config", "Request method is required for HTTP server type.")
		}

		if data.HTTPConfig.AuthenticationType == nil || *data.HTTPConfig.AuthenticationType == "" {
			resp.Diagnostics.AddError("authentication_type is missing from http_config", "Authentication type is required for HTTP server type.")
		}

		if data.HTTPConfig.Username == nil || *data.HTTPConfig.Username == "" {
			resp.Diagnostics.AddError("username is missing from http_config", "Username is required for HTTP server type.")
		}

		if data.HTTPConfig.Password == nil || *data.HTTPConfig.Password == "" {
			resp.Diagnostics.AddError("password is missing from http_config", "Password is required for HTTP server type.")
		}

		if data.HTTPConfig.ApiKey == nil || *data.HTTPConfig.ApiKey == "" {
			resp.Diagnostics.AddError("api_key is missing from http_config", "API key is required for HTTP server type.")
		}

		if data.HTTPConfig.TransportMode == nil || *data.HTTPConfig.TransportMode == "" {
			resp.Diagnostics.AddError("transport_mode is missing from http_config", "Transport mode is required for HTTP server type.")
		}

		if data.HTTPConfig.ActionBody == nil || *data.HTTPConfig.ActionBody == "" {
			resp.Diagnostics.AddError("action_body is missing from http_config", "Action body is required for HTTP server type.")
		}
	}
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
