package provider

import (
	"context"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ApiScopeResourceSetup() resource.Resource {
	return &ApiScopeResource{}
}

func (r *ApiScopeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		providerInput = req.ProviderData.(*ElvidProviderInput)
	}
}

func (r *ApiScopeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_apiscope"
}

func (r *ApiScopeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "The name of the API scope",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "A description of what this API scope is used for. Please include information about what it gives access to, and in what way it differs from similar API scopes, if any.",
			},
			"user_claims": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // This means empty string set as default
				Description: "User claims that are included in the token when logging in with a machine/user client that has this API scope. (The token will include the super set of claims from all granted scopes).",
			},
			"allow_machine_clients": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the API scope is intended for machine clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			"allow_user_clients": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the API scope is intended for user clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			"resource_taint_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Default:     stringdefault.StaticString("1"),
				Description: "A change in value for this field will force recreating the resource",
			},
		},
	}
}

// ModifyPlan runs ElvID's server-side H2 validation against the planned values
// at plan time, so naming-convention and AD-group-existence errors surface in
// the PR's speculative plan instead of mid-apply on merge. We only validate on
// create (state is null); updates and deletes go through their own checks at
// apply. See ADR 2026-05-CORE-2650 (H2) in the elvid repo.
func (r *ApiScopeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Plan removal (delete) — nothing to validate.
	if req.Plan.Raw.IsNull() {
		return
	}
	// Update — H2 only applies on create; server's validate endpoint no-ops if
	// the scope name already exists, but skip the call entirely to keep plans fast
	// and offline-resilient for the common case.
	if !req.State.Raw.IsNull() {
		return
	}
	// Provider not yet configured (e.g. terraform validate without auth) — skip
	// the call rather than hard-fail; the apply will still enforce.
	if providerInput == nil || providerInput.AccessTokenAD == "" {
		return
	}

	var plan ApiScopeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dto := plan.DtoFromApiScopeResource()
	if err := elvidapiclient.ValidateNewApiScope(ctx, providerInput.ElvIDAuthority, providerInput.AccessTokenAD, dto); err != nil {
		resp.Diagnostics.AddError("ApiScope plan-time validation failed", err.Error())
	}
}

func (r *ApiScopeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApiScopeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiScopeRequestDto := plan.DtoFromApiScopeResource()
	_, err := elvidapiclient.CreateOrUpdateApiScope(ctx, providerInput.ElvIDAuthority, providerInput.AccessTokenAD, apiScopeRequestDto)
	if err != nil {
		resp.Diagnostics.AddError("Creating ApiScope resulted in an error", err.Error())
		return
	}

	plan.Id = types.StringValue(plan.Name.ValueString())
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ApiScopeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApiScopeResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiScopeRequestDto := plan.DtoFromApiScopeResource()
	_, err := elvidapiclient.CreateOrUpdateApiScope(ctx, providerInput.ElvIDAuthority, providerInput.AccessTokenAD, apiScopeRequestDto)
	if err != nil {
		resp.Diagnostics.AddError("Update ApiScope resulted in an error", err.Error())
		return
	}

	plan.Id = types.StringValue(plan.Name.ValueString())
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

	apiScopeResponseDto, err := elvidapiclient.ReadApiScope(ctx, providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading ApiScope resulted in an error", err.Error())
		return
	}

	if apiScopeResponseDto == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.FromElvidApiClient(apiScopeResponseDto)
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

	err := elvidapiclient.DeleteApiScope(ctx, providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
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
