package email_configuration_resource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/kibblator/terraform-provider-ory/internal/provider/custom_validators"

	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"
	"github.com/ory/client-go"
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
	SenderName    types.String `tfsdk:"sender_name"`
	SenderAddress types.String `tfsdk:"sender_address"`
	Host          types.String `tfsdk:"host"`
	Port          types.String `tfsdk:"port"`
	Security      types.String `tfsdk:"security"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
}

type HTTPConfig struct {
	Url                types.String `tfsdk:"url"`
	RequestMethod      types.String `tfsdk:"request_method"`
	AuthenticationType types.String `tfsdk:"authentication_type"`
	BasicAuth          *BasicAuth   `tfsdk:"basic_auth"`
	ApiKey             *APIKey      `tfsdk:"api_key"`
	ActionBody         types.String `tfsdk:"action_body"`
}

type APIKey struct {
	TransportMode types.String `tfsdk:"transport_mode"`
	Name          types.String `tfsdk:"name"`
	Value         types.String `tfsdk:"value"`
}

type BasicAuth struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type SMTPHeader struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type emailConfigurationResourceModel struct {
	ID          types.String  `tfsdk:"id"`
	LastUpdated types.String  `tfsdk:"last_updated"`
	ServerType  types.String  `tfsdk:"server_type"`
	SMTPConfig  *SMTPConfig   `tfsdk:"smtp_config"`
	HTTPConfig  *HTTPConfig   `tfsdk:"http_config"`
	SMTPHeaders *[]SMTPHeader `tfsdk:"smtp_headers"`
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
			"id": schema.StringAttribute{
				Description: "String identifier of the email configuration resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the email configuration settings.",
				Computed:    true,
			},
			"server_type": schema.StringAttribute{
				Description: "The type of the email server.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("smtp", "http", "default"),
				},
			},
			"smtp_config": schema.SingleNestedAttribute{
				Description: "SMTP configuration block (optional, but fields required if present).",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sender_name": schema.StringAttribute{
						Description: "The name of the sender.",
						Optional:    true,
						Computed:    true,
					},
					"sender_address": schema.StringAttribute{
						Description: "The email address of the sender.",
						Optional:    true,
						Computed:    true,
					},
					"host": schema.StringAttribute{
						Description: "The SMTP server host.",
						Optional:    true,
						Computed:    true,
					},
					"port": schema.StringAttribute{
						Description: "The SMTP server port.",
						Optional:    true,
						Computed:    true,
					},
					"security": schema.StringAttribute{
						Description: "The security type of the SMTP server.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("starttls", "starttls_notrust", "cleartext", "implicittls", "implicittls_notrust"),
						},
					},
					"username": schema.StringAttribute{
						Description: "The username for the SMTP server.",
						Optional:    true,
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "The password for the SMTP server.",
						Optional:    true,
						Sensitive:   true,
						Computed:    true,
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
						Computed:    true,
					},
					"request_method": schema.StringAttribute{
						Description: "The request method for the HTTP server.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("GET", "POST", "PUT", "PATCH"),
						},
					},
					"authentication_type": schema.StringAttribute{
						Description: "The authentication type for the HTTP server.",
						Optional:    true,
						Computed:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "basic_auth", "api_key"),
						},
					},
					"api_key": schema.SingleNestedAttribute{
						Description: "The API key for the HTTP server.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"transport_mode": schema.StringAttribute{
								Description: "The transport mode for the HTTP server.",
								Required:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("header", "cookie"),
								},
							},
							"name": schema.StringAttribute{
								Description: "The name of the API Key.",
								Required:    true,
							},
							"value": schema.StringAttribute{
								Description: "The value of the API Key.",
								Required:    true,
								Sensitive:   true,
							},
						},
					},
					"basic_auth": schema.SingleNestedAttribute{
						Description: "The basic auth configuration for the HTTP server.",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description: "The username for the HTTP server auth.",
								Required:    true,
							},
							"password": schema.StringAttribute{
								Description: "The password for the HTTP server auth.",
								Required:    true,
								Sensitive:   true,
							},
						},
					},
					"action_body": schema.StringAttribute{
						Description: "The base64 encoded action body for the HTTP server.",
						Optional:    true,
						Computed:    true,
						Validators:  []validator.String{custom_validators.Base64Validator{}},
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

	if data.ServerType.ValueString() == "default" {
		if data.SMTPConfig != nil || data.HTTPConfig != nil {
			resp.Diagnostics.AddError("smtp_config and http_config are not allowed with server_type default", "SMTP and HTTP configurations are not allowed with default server type.")
		}
	}

	if data.ServerType.ValueString() == "smtp" {
		if data.HTTPConfig != nil {
			resp.Diagnostics.AddError("http_config is not allowed with server_type smtp", "HTTP configuration is not allowed with SMTP server type.")
			return
		}

		if data.SMTPConfig == nil {
			resp.Diagnostics.AddError("smtp_config is missing", "SMTP configuration is required with SMTP server type.")
			return
		}

		if data.SMTPConfig.SenderName.ValueString() == "" {
			resp.Diagnostics.AddError("sender_name is missing from smtp_config", "Sender name is required for SMTP server type.")
		}

		if data.SMTPConfig.SenderAddress.ValueString() == "" {
			resp.Diagnostics.AddError("sender_address is missing from smtp_config", "Sender address is required for SMTP server type.")
		}

		if data.SMTPConfig.Host.ValueString() == "" {
			resp.Diagnostics.AddError("host is missing from smtp_config", "Host is required for SMTP server type.")
		}

		if data.SMTPConfig.Port.ValueString() == "" {
			resp.Diagnostics.AddError("port is missing from smtp_config", "Port is required for SMTP server type.")
		}

		if data.SMTPConfig.Security.ValueString() == "" {
			resp.Diagnostics.AddError("security is missing from smtp_config", "Security is required for SMTP server type.")
		}

		if data.SMTPConfig.Username.ValueString() == "" {
			resp.Diagnostics.AddError("username is missing from smtp_config", "Username is required for SMTP server type.")
		}

		if data.SMTPConfig.Password.ValueString() == "" {
			resp.Diagnostics.AddError("password is missing from smtp_config", "Password is required for SMTP server type.")
		}
	}

	if data.ServerType.ValueString() == "http" {
		if data.SMTPConfig != nil {
			resp.Diagnostics.AddError("smtp_config is not allowed with server_type http", "SMTP configuration is not allowed with HTTP server type.")
			return
		}

		if data.HTTPConfig == nil {
			resp.Diagnostics.AddError("http_config is missing", "HTTP configuration is required with HTTP server type.")
			return
		}

		if data.HTTPConfig.Url.ValueString() == "" {
			resp.Diagnostics.AddError("url is missing from http_config", "URL is required for HTTP server type.")
		}

		if data.HTTPConfig.RequestMethod.ValueString() == "" {
			resp.Diagnostics.AddError("request_method is missing from http_config", "Request method is required for HTTP server type.")
		}

		if data.HTTPConfig.AuthenticationType.ValueString() == "" {
			resp.Diagnostics.AddError("authentication_type is missing from http_config", "Authentication type is required for HTTP server type.")
		}

		if data.HTTPConfig.AuthenticationType.ValueString() == "none" {
			if data.HTTPConfig.ApiKey != nil || data.HTTPConfig.BasicAuth != nil {
				resp.Diagnostics.AddError("authentication_method not allowed", "cannot specify an authentication method if authentication_type is set to 'none'")
			}
		}

		if data.HTTPConfig.AuthenticationType.ValueString() == "basic_auth" {
			if data.HTTPConfig.ApiKey != nil {
				resp.Diagnostics.AddError("api_key not allowed", "cannot specify api_key block when authentication_type is 'basic'")
			}
		}

		if data.HTTPConfig.AuthenticationType.ValueString() == "api_key" {
			if data.HTTPConfig.BasicAuth != nil {
				resp.Diagnostics.AddError("basic_auth not allowed", "cannot specify basic_auth block when authentication_type is 'api_key'")
			}
		}

		if data.HTTPConfig.ActionBody.ValueString() == "" {
			resp.Diagnostics.AddError("action_body is missing from http_config", "Action body is required for HTTP server type.")
		}
	}
}

// Create implements resource.Resource.
func (r *emailConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan emailConfigurationResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var patch []client.JsonPatch

	var oryConfig orytypes.Config
	orytypes.TransformToConfig(r.oryClient.ProjectConfig.Services.Identity.Config, &oryConfig)

	if plan.ServerType.ValueString() == "default" {
		if oryConfig.Courier.SMTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/smtp",
			})
		}

		if oryConfig.Courier.HTTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/http",
			})
		}

		if oryConfig.Courier.DeliveryStrategy != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/delivery_strategy",
			})
		}
	}

	if plan.ServerType.ValueString() == "smtp" {
		if oryConfig.Courier.HTTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/http",
			})
		}

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/delivery_strategy",
			Value: "smtp",
		})

		var smtpConfig orytypes.SMTP
		SmtpConfigToApi(plan, &smtpConfig)

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/smtp",
			Value: smtpConfig,
		})
	}

	if plan.ServerType.ValueString() == "http" {
		if oryConfig.Courier.SMTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/smtp",
			})
		}

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/delivery_strategy",
			Value: "http",
		})

		headersMap := make(map[string]string)

		for _, header := range *plan.SMTPHeaders {
			headersMap[header.Key.ValueString()] = header.Value.ValueString()
		}

		var httpConfig orytypes.HTTP
		HttpConfigToApi(plan, &httpConfig)

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/http",
			Value: httpConfig,
		})
	}

	// If no changes detected, exit early
	if len(patch) == 0 {
		resp.Diagnostics.AddWarning(
			"No Changes Detected",
			"Update was triggered but no changes were detected between the plan and the current state.",
		)
		return
	}

	projectUpdate, _, err := r.oryClient.APIClient.ProjectAPI.PatchProjectWithRevision(ctx, r.oryClient.ProjectID, r.oryClient.ProjectConfig.RevisionId).JsonPatch(patch).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ory email configuration",
			"Could not update ory email configuration, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue("email_configuration_settings")
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	var updatedConfig orytypes.Config
	orytypes.TransformToConfig(projectUpdate.Project.Services.Identity.Config, &updatedConfig)

	if updatedConfig.Courier.DeliveryStrategy == nil {
		plan.ServerType = types.StringValue("default")
	} else {
		plan.ServerType = types.StringValue(*updatedConfig.Courier.DeliveryStrategy)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read implements resource.Resource.
func (r *emailConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading email configuration resource")

	// Retrieve current state
	var state emailConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch current project configuration from ORY
	project, _, err := r.oryClient.APIClient.ProjectAPI.GetProject(ctx, r.oryClient.ProjectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching ORY registration config",
			"Could not retrieve ORY email configuration: "+err.Error(),
		)
		return
	}

	var projectConfig orytypes.Config
	orytypes.TransformToConfig(project.Services.Identity.Config, &projectConfig)

	// Update the state with current configuration values
	state.ID = types.StringValue("email_configuration_settings")
	serverType := "default"

	if projectConfig.Courier.DeliveryStrategy != nil {
		serverType = *projectConfig.Courier.DeliveryStrategy
	}

	state.ServerType = types.StringValue(serverType)

	if serverType == "smtp" {
		err := ApiToSmtpConfig(projectConfig.Courier.SMTP, &state)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing SMTP URL",
				"Could not parse SMTP URL: "+err.Error(),
			)
			return
		}
	}

	if serverType == "http" {
		ApiToHttpConfig(projectConfig.Courier.HTTP, &state)
	}

	tflog.Debug(ctx, "Updated State", map[string]interface{}{
		"state": state,
	})

	// Set the updated state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update implements resource.Resource.
func (r *emailConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan emailConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var patch []client.JsonPatch
	var oryConfig orytypes.Config
	orytypes.TransformToConfig(r.oryClient.ProjectConfig.Services.Identity.Config, &oryConfig)

	if plan.ServerType.ValueString() == "default" {
		if oryConfig.Courier.SMTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/smtp",
			})
		}

		if oryConfig.Courier.HTTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/http",
			})
		}

		if oryConfig.Courier.DeliveryStrategy != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/delivery_strategy",
			})
		}
	}

	if plan.ServerType.ValueString() == "smtp" {
		if oryConfig.Courier.HTTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/http",
			})
		}

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/delivery_strategy",
			Value: "smtp",
		})

		var smtpConfig orytypes.SMTP
		SmtpConfigToApi(plan, &smtpConfig)

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/smtp",
			Value: smtpConfig,
		})
	}

	if plan.ServerType.ValueString() == "http" {
		if oryConfig.Courier.SMTP != nil {
			patch = append(patch, client.JsonPatch{
				Op:   "remove",
				Path: "/services/identity/config/courier/smtp",
			})
		}

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/delivery_strategy",
			Value: "http",
		})

		var httpConfig orytypes.HTTP
		HttpConfigToApi(plan, &httpConfig)

		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/courier/http",
			Value: httpConfig,
		})
	}

	_, _, err := r.oryClient.APIClient.ProjectAPI.PatchProjectWithRevision(ctx, r.oryClient.ProjectID, r.oryClient.ProjectConfig.RevisionId).JsonPatch(patch).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ory email configuration",
			"Could not update ory email configuration, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue("email_configuration_settings")
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (r *emailConfigurationResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
}

// ImportState implements resource.ResourceWithImportState.
func (r *emailConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
