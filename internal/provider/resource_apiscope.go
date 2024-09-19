package provider

import (
	"context"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ApiScopeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_apiscope"
}

func (r *ApiScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				ForceNew:    true,
				Description: "The name of the API scope",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of what this API scope is used for. Please include information about what it gives access to, and in what way it differs from similar API scopes, if any.",
			},
			"user_claims": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "User claims that are included in the token when logging in with a machine/user client that has this API scope. (The token will include the superset of claims from all granted scopes).",
			},
			"allow_machine_clients": schema.BoolAttribute{
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for machine clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			"allow_user_clients": schema.BoolAttribute{
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for user clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			"resource_taint_version": schema.StringAttribute{
				Optional:    true,
				ForceNew:    true,
				Default:     "1",
				Description: "A change in value for this field will force recreating the resource",
			},
		},
	}
}

func (r *ApiScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiScopeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiScopeDto := plan.DtoFromApiScopeResource()
	apiScope, err := elvidapiclient.CreateOrUpdateApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, apiScopeDto)
	if err != nil {
		resp.Diagnostics.AddError("Creating ApiScope resulted in an error", err.Error())
		return
	}

	plan.Id = types.StringValue(apiScope.Name)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ApiScopeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApiScopeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiScope, diag := elvidapiclient.ReadApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if diag.HasError() {
		return
	}

	if apiScope == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.FromElvidApiClient(apiScope)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ApiScopeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ApiScopeResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := elvidapiclient.DeleteApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Deleting API scope resulted in an error", err.Error())
	}
}

type ApiScopeResource struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	UserClaims           types.Set    `tfsdk:"user_claims"`
	AllowMachineClients  types.Bool   `tfsdk:"allow_machine_clients"`
	AllowUserClients     types.Bool   `tfsdk:"allow_user_clients"`
	ResourceTaintVersion types.String `tfsdk:"resource_taint_version"`
}

func (as *ApiScopeResource) DtoFromApiScopeResource() *elvidapiclient.ApiScopeDto {
	return &elvidapiclient.ApiScopeDto{
		Name:                as.Name.ValueString(),
		Description:         as.Description.ValueString(),
		UserClaims:          convertSetToStringArray(as.UserClaims),
		AllowMachineClients: as.AllowMachineClients.ValueBool(),
		AllowUserClients:    as.AllowUserClients.ValueBool(),
	}
}

func (as *ApiScopeResource) FromElvidApiClient(apiScope *elvidapiclient.ApiScopeDto) {
	as.Id = types.StringValue(apiScope.Name)
	as.Name = types.StringValue(apiScope.Name)
	as.Description = types.StringValue(apiScope.Description)
	as.UserClaims = convertStringArrayToSet(apiScope.UserClaims)
	as.AllowMachineClients = types.BoolValue(apiScope.AllowMachineClients)
	as.AllowUserClients = types.BoolValue(apiScope.AllowUserClients)
}
