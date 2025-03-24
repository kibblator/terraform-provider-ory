package registration_resource

import (
	"context"
	"fmt"
	"strconv"
	"time"

	orytypes "github.com/kibblator/terraform-provider-ory/internal/provider/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ory/client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &registrationResource{}
	_ resource.ResourceWithConfigure   = &registrationResource{}
	_ resource.ResourceWithImportState = &registrationResource{}
)

// NewRegistrationResource is a helper function to simplify the provider implementation.
func NewRegistrationResource() resource.Resource {
	return &registrationResource{}
}

// registrationResource is the resource implementation.
type registrationResource struct {
	oryClient *orytypes.OryClient
}

// registrationResourceModel maps the resource schema data.
type registrationResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	EnableRegistration  types.Bool   `tfsdk:"enable_registration"`
	EnablePasswordAuth  types.Bool   `tfsdk:"enable_password_auth"`
	EnablePostSigninReg types.Bool   `tfsdk:"enable_post_signin_reg"`
	EnableLoginHints    types.Bool   `tfsdk:"enable_login_hints"`
}

// Configure adds the provider configured client to the resource.
func (r *registrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *registrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration"
}

// Schema defines the schema for the resource.
func (r *registrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "String identifier of the registration resource.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Timestamp of the last Terraform update of the registration settings.",
				Computed:    true,
			},
			"enable_registration": schema.BoolAttribute{
				Description: "If enabled, users can sign up using the selfservice UIs.",
				Optional:    true,
				Computed:    true,
			},
			"enable_password_auth": schema.BoolAttribute{
				Description: "If enabled, users will be able to sign in and register using a password.",
				Optional:    true,
				Computed:    true,
			},
			"enable_post_signin_reg": schema.BoolAttribute{
				Description: "If enabled, users will be automatically logged in after they register.",
				Optional:    true,
				Computed:    true,
			},
			"enable_login_hints": schema.BoolAttribute{
				Description: "Login hints provide additional information to users when they try to sign up with an account, that already exists",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func findHookIndex(hooks []orytypes.Hook, target string) int {
	for i, hook := range hooks {
		if hook.Hook == target {
			return i
		}
	}
	return -1
}

// Create a new resource.
func (r *registrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan registrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var patch []client.JsonPatch

	// Conditionally append patches only if values are explicitly set
	if !plan.EnableRegistration.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/flows/registration/enabled",
			Value: plan.EnableRegistration.ValueBool(),
		})
	}

	if !plan.EnableLoginHints.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/flows/registration/login_hints",
			Value: plan.EnableLoginHints.ValueBool(),
		})
	}

	if !plan.EnablePasswordAuth.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/methods/password/enabled",
			Value: plan.EnablePasswordAuth.ValueBool(),
		})
	}

	var oryConfig orytypes.Config
	orytypes.TransformToConfig(r.oryClient.ProjectConfig.Services.Identity.Config, &oryConfig)

	if !plan.EnablePostSigninReg.IsNull() {
		if plan.EnablePostSigninReg.ValueBool() && findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session") == -1 {
			patch = append(patch, client.JsonPatch{
				Op:   "add",
				Path: "/services/identity/config/selfservice/flows/registration/after/password/hooks/0",
				Value: orytypes.Hook{
					Hook: "session",
				},
			})

		} else {
			index := findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session")
			if index != -1 {
				patch = append(patch, client.JsonPatch{
					Op:   "remove",
					Path: "/services/identity/config/selfservice/flows/registration/after/password/hooks/" + strconv.Itoa(index),
				})
			}
		}
	}

	tflog.Debug(ctx, "Generated Patch", map[string]interface{}{
		"patch": patch,
	})

	projectUpdate, _, err := r.oryClient.APIClient.ProjectAPI.PatchProjectWithRevision(ctx, r.oryClient.ProjectID, r.oryClient.ProjectConfig.RevisionId).JsonPatch(patch).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ory registration config",
			"Could not update ory registration config, unexpected error: "+err.Error(),
		)
		return
	}

	orytypes.TransformToConfig(projectUpdate.Project.Services.Identity.Config, &oryConfig)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ory registration config",
			"Could not update ory registration config, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue("registration_settings")

	enablePostSigninReg := findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session") != -1

	plan.EnableLoginHints = types.BoolValue(oryConfig.SelfService.Flows.Registration.LoginHints)
	plan.EnableRegistration = types.BoolValue(oryConfig.SelfService.Flows.Registration.Enabled)
	plan.EnablePasswordAuth = types.BoolValue(oryConfig.SelfService.Methods.Password.Enabled)
	plan.EnablePostSigninReg = types.BoolValue(enablePostSigninReg)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *registrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Reading registration resource")

	// Retrieve current state
	var state registrationResourceModel
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
			"Could not retrieve ORY registration configuration: "+err.Error(),
		)
		return
	}

	var oryConfig orytypes.Config
	orytypes.TransformToConfig(project.Services.Identity.Config, &oryConfig)

	enablePostSigninReg := findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session") != -1

	// Update the state with current configuration values
	state.ID = types.StringValue("registration_settings")

	state.EnableRegistration = types.BoolValue(oryConfig.SelfService.Flows.Registration.Enabled)
	state.EnableLoginHints = types.BoolValue(oryConfig.SelfService.Flows.Registration.LoginHints)
	state.EnablePasswordAuth = types.BoolValue(oryConfig.SelfService.Methods.Password.Enabled)
	state.EnablePostSigninReg = types.BoolValue(enablePostSigninReg)

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

func (r *registrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the desired state from the plan
	var plan registrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize the patch list
	var patch []client.JsonPatch

	tflog.Debug(ctx, "Update Plan", map[string]interface{}{
		"plan": plan,
	})

	// Compare the plan with the current state and add patches for changes
	if !plan.EnableRegistration.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/flows/registration/enabled",
			Value: plan.EnableRegistration.ValueBool(),
		})
	}

	if !plan.EnableLoginHints.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/flows/registration/login_hints",
			Value: plan.EnableLoginHints.ValueBool(),
		})
	}

	if !plan.EnablePasswordAuth.IsNull() {
		patch = append(patch, client.JsonPatch{
			Op:    "replace",
			Path:  "/services/identity/config/selfservice/methods/password/enabled",
			Value: plan.EnablePasswordAuth.ValueBool(),
		})
	}

	tflog.Debug(ctx, "Generated Patch", map[string]interface{}{
		"patch": patch,
	})

	var oryConfig orytypes.Config
	orytypes.TransformToConfig(r.oryClient.ProjectConfig.Services.Identity.Config, &oryConfig)

	if !plan.EnablePostSigninReg.IsNull() {
		if plan.EnablePostSigninReg.ValueBool() && findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session") == -1 {
			// Add the "session" hook if it doesn't exist
			patch = append(patch, client.JsonPatch{
				Op:   "add",
				Path: "/services/identity/config/selfservice/flows/registration/after/password/hooks/0",
				Value: orytypes.Hook{
					Hook: "session",
				},
			})
		} else {
			// Remove the "session" hook if it exists
			index := findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session")
			if index != -1 {
				patch = append(patch, client.JsonPatch{
					Op:   "remove",
					Path: "/services/identity/config/selfservice/flows/registration/after/password/hooks/" + strconv.Itoa(index),
				})
			}
		}
	}

	// If no changes detected, exit early
	if len(patch) == 0 {
		resp.Diagnostics.AddWarning(
			"No Changes Detected",
			"Update was triggered but no changes were detected between the plan and the current state.",
		)
		return
	}

	//get latest revision
	project, _, _ := r.oryClient.APIClient.ProjectAPI.GetProject(ctx, r.oryClient.ProjectID).Execute()
	projectUpdate, _, err := r.oryClient.APIClient.ProjectAPI.PatchProjectWithRevision(ctx, r.oryClient.ProjectID, project.RevisionId).JsonPatch(patch).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ory registration config",
			"Could not update ory registration config, unexpected error: "+err.Error(),
		)
		return
	}

	orytypes.TransformToConfig(projectUpdate.Project.Services.Identity.Config, &oryConfig)

	enablePostSigninReg := findHookIndex(oryConfig.SelfService.Flows.Registration.After.Password.Hooks, "session") != -1

	// Update plan with the extracted values
	plan.EnableLoginHints = types.BoolValue(oryConfig.SelfService.Flows.Registration.LoginHints)
	plan.EnableRegistration = types.BoolValue(oryConfig.SelfService.Flows.Registration.Enabled)
	plan.EnablePasswordAuth = types.BoolValue(oryConfig.SelfService.Methods.Password.Enabled)
	plan.EnablePostSigninReg = types.BoolValue(enablePostSigninReg)

	// Update ID and LastUpdated
	plan.ID = types.StringValue("registration_settings")
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set the updated plan to the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *registrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *registrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
