package provider

import (
	"context"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var providerInput *ElvidProviderInput

func NewClientSecretResource() resource.Resource {
	return &ClientSecretResource{}
}

func (r *ClientSecretResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_clientsecret"
}

func (r *ClientSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		providerInput = req.ProviderData.(*ElvidProviderInput)
	}
}

func (r *ClientSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Required:    true,
				Description: "The (entity) id of the client. Note that this should be set to client.id and not client.client_id",
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
			"secret_value": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Sensitive: true,
			},
			"hashed_value_starts_with": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ClientSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClientSecretResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdClientSecret, err := elvidapiclient.CreateClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, plan.ClientID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Creating clientsecret resulted in an error", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdClientSecret.Id))
	plan.SecretValue = types.StringValue(createdClientSecret.Value)
	plan.HashedValueStartsWith = types.StringValue(createdClientSecret.HashedValueStartsWith)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ClientSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClientSecretResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientSecret, err := elvidapiclient.ReadClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.ClientID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading clientsecret resulted in an error", err.Error())
		return
	}

	if clientSecret == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	if clientSecret.Value == "" {
		clientSecret.Value = state.SecretValue.ValueString()
	}

	if providerInput.RunHashedSecretValidation && clientSecret.HashedValueStartsWith != state.HashedValueStartsWith.ValueString() {
		resp.Diagnostics.AddError("HashedValueStartsWith has changed. Recreate secret with setting run_hashed_secret_validation = false in provider config", "")
		return
	}

	state.FromElvidApiClientOutput(clientSecret)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *ClientSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	panic("unimplemented")
}

func (r *ClientSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClientSecretResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := elvidapiclient.DeleteClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.ClientID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Deleting clientsecret resulted in an error", err.Error())
	}
}

type ClientSecretResource struct {
	ID                    types.String `tfsdk:"id"`
	ClientID              types.String `tfsdk:"client_id"`
	ResourceTaintVersion  types.String `tfsdk:"resource_taint_version"`
	SecretValue           types.String `tfsdk:"secret_value"`
	HashedValueStartsWith types.String `tfsdk:"hashed_value_starts_with"`
}

func (cs *ClientSecretResource) FromElvidApiClientOutput(clientSecret *elvidapiclient.ClientSecret) {
	cs.HashedValueStartsWith = types.StringValue(clientSecret.HashedValueStartsWith)
	cs.SecretValue = types.StringValue(clientSecret.Value)
}
