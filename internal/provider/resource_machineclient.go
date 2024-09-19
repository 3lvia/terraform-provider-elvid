package provider

import (
	"context"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewMachineClientResource() resource.Resource {
	return &MachineClientResource{}
}

func (r *MachineClientResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_machineclient"
}

func (r *MachineClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		providerInput = req.ProviderData.(*ElvidProviderInput)
	}
}

func (r *MachineClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the client",
			},
			"test_user_login_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When this is enabled it's possible to mechanically login with a test user to this client. This is done by using grant-type password for the token endpoint. See ElvID space on confluence for details",
			},
			"is_delegation_client": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When this is enabled the client can only use the delegation grant type. This is used when a already logged inn user will create a long-lived delegation access_token",
			},
			"access_token_life_time": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3600),
				Description: "Number of seconds before an access token expires. Default and maximum is 3600 seconds (1 hour) for regular machine clients.",
			},
			"scopes": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
			},
			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "The client_id of the client, used during client_credentials auth. Note this is different from the (entity) id of the client",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"token_endpoint": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"client_claims": schema.SetNestedBlock{
				Description: "Client claims that will always be added in the access_token. The claim type must start with 'client_'.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Type of the claim, must start with 'client_'.",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "List of values associated with the claim.",
						},
					},
				},
			},
		},
	}
}

func (r *MachineClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MachineClientResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineClientRequestDto := plan.DtoFromMachineClientResource(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	machineClientResponseDto, err := elvidapiclient.CreateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientRequestDto)
	if err != nil {
		resp.Diagnostics.AddError("Creating machineclient resulted in an error", err.Error())
		return
	}

	plan.Id = types.StringValue(strconv.Itoa(machineClientResponseDto.Id))
	plan.ClientID = types.StringValue(machineClientResponseDto.ClientId)
	plan.TokenEndpoint = types.StringValue(providerInput.ElvIDAuthority + "/connect/token")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MachineClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MachineClientResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineClientResponseDto, err := elvidapiclient.ReadMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading machineclient resulted in an error", err.Error())
		return
	}

	if machineClientResponseDto == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.MachineClientResourceFromDto(machineClientResponseDto)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *MachineClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MachineClientResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	machineClientRequestDto := plan.DtoFromMachineClientResource(ctx, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	machineClientRequestDto.Id, _ = strconv.Atoi(plan.Id.ValueString())

	_, err := elvidapiclient.UpdateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientRequestDto)
	if err != nil {
		resp.Diagnostics.AddError("Updating machineclient resulted in an error", err.Error())
		return
	}

	plan.TokenEndpoint = types.StringValue(providerInput.ElvIDAuthority + "/connect/token")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *MachineClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MachineClientResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := elvidapiclient.DeleteMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Deleting machineclient resulted in an error", err.Error())
	}
}

type MachineClientResource struct {
	Id                   types.String          `tfsdk:"id"`
	Name                 types.String          `tfsdk:"name"`
	TestUserLoginEnabled types.Bool            `tfsdk:"test_user_login_enabled"`
	IsDelegationClient   types.Bool            `tfsdk:"is_delegation_client"`
	AccessTokenLifeTime  types.Int64           `tfsdk:"access_token_life_time"`
	Scopes               types.Set             `tfsdk:"scopes"`
	ClientID             types.String          `tfsdk:"client_id"`
	ResourceTaintVersion types.String          `tfsdk:"resource_taint_version"`
	TokenEndpoint        types.String          `tfsdk:"token_endpoint"`
	ClientClaims         []ClientClaimResource `tfsdk:"client_claims"`
}

type ClientClaimResource struct {
	Type   types.String   `tfsdk:"type"`
	Values []types.String `tfsdk:"values"`
}

func (mc *MachineClientResource) DtoFromMachineClientResource(ctx context.Context, diagnostics diag.Diagnostics) *elvidapiclient.MachineClientDto {
	return &elvidapiclient.MachineClientDto{
		ClientName:           mc.Name.ValueString(),
		Scopes:               convertSetToStringList(mc.Scopes),
		TestUserLoginEnabled: mc.TestUserLoginEnabled.ValueBool(),
		AccessTokenLifeTime:  int(mc.AccessTokenLifeTime.ValueInt64()),
		IsDelegationClient:   mc.IsDelegationClient.ValueBool(),
		ClientClaims:         DtoFromClientClaimResource(mc.ClientClaims, diagnostics),
	}
}

func (mc *MachineClientResource) MachineClientResourceFromDto(client *elvidapiclient.MachineClientDto) {
	mc.Name = types.StringValue(client.ClientName)
	mc.TestUserLoginEnabled = types.BoolValue(client.TestUserLoginEnabled)
	mc.AccessTokenLifeTime = types.Int64Value(int64(client.AccessTokenLifeTime))
	mc.IsDelegationClient = types.BoolValue(client.IsDelegationClient)
	mc.ClientID = types.StringValue(client.ClientId)
}

func convertSetToStringList(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}

func DtoFromClientClaimResource(clientClaims []ClientClaimResource, diagnostics diag.Diagnostics) []elvidapiclient.ClientClaimDto {
	var clientClaimDtos []elvidapiclient.ClientClaimDto

	if len(clientClaims) == 0 {
		return []elvidapiclient.ClientClaimDto{}
	}

	for _, clientClaim := range clientClaims {
		// Validate and process clientClaimResource.Type
		if clientClaim.Type.IsNull() || clientClaim.Type.IsUnknown() {
			diagnostics.AddError(
				"Invalid Client Claim Type",
				"Client claim type must be provided and cannot be null or unknown.",
			)
			continue
		}

		// Initialize a slice to hold the values
		var values []string

		// Iterate over each value in clientClaimResource.Values
		for _, v := range clientClaim.Values {
			if v.IsNull() || v.IsUnknown() {
				diagnostics.AddWarning(
					"Unknown or Null Value",
					"One of the values in clientClaimResource.Values is null or unknown and will be skipped.",
				)
				continue
			}
			values = append(values, v.ValueString())
		}

		// Empty list instead of nil for no elements
		if values == nil {
			values = []string{}
		}

		// Create a ClientClaim object
		clientClaimDto := elvidapiclient.ClientClaimDto{
			Type:   clientClaim.Type.ValueString(),
			Values: values,
		}

		// Append the ClientClaim to the slice
		clientClaimDtos = append(clientClaimDtos, clientClaimDto)
	}

	return clientClaimDtos
}
