package provider

import (
	"context"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var providerInput *ElvidProviderInput

func NewMachineClientResource() resource.Resource {
	return &MachineClientResource{}
}

func (r *MachineClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		providerInput = req.ProviderData.(*ElvidProviderInput)
	}
}

func (r *MachineClientResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_machineclient"
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
			"client_claims": schema.SetNestedAttribute{
				Optional:    true,
				Description: "Client claims is claims for the client that will always be added in the access_token. The claim type must start with 'client_'.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
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

	machineClientInput := plan.ToElvidApiClient()
	machineClient, err := elvidapiclient.CreateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)
	if err != nil {
		resp.Diagnostics.AddError("Creating machineclient resulted in an error", err.Error())
		return
	}

	plan.Id = types.StringValue(strconv.Itoa(machineClient.Id))
	plan.ClientID = types.StringValue(machineClient.ClientId)
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

	machineClient, err := elvidapiclient.ReadMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading machineclient resulted in an error", err.Error())
		return
	}

	if machineClient == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.FromElvidApiClient(machineClient)
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

	machineClientInput := plan.ToElvidApiClient()
	machineClientInput.Id, _ = strconv.Atoi(plan.Id.ValueString())

	_, err := elvidapiclient.UpdateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)
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
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	TestUserLoginEnabled types.Bool   `tfsdk:"test_user_login_enabled"`
	IsDelegationClient   types.Bool   `tfsdk:"is_delegation_client"`
	AccessTokenLifeTime  types.Int64  `tfsdk:"access_token_life_time"`
	Scopes               types.Set    `tfsdk:"scopes"`
	ClientID             types.String `tfsdk:"client_id"`
	ResourceTaintVersion types.String `tfsdk:"resource_taint_version"`
	TokenEndpoint        types.String `tfsdk:"token_endpoint"`
	ClientClaims         types.Set    `tfsdk:"client_claims"`
}

func (mc *MachineClientResource) ToElvidApiClient() *elvidapiclient.MachineClient {
	return &elvidapiclient.MachineClient{
		ClientName:           mc.Name.ValueString(),
		Scopes:               convertSetToStringSlice(mc.Scopes),
		TestUserLoginEnabled: mc.TestUserLoginEnabled.ValueBool(),
		AccessTokenLifeTime:  int(mc.AccessTokenLifeTime.ValueInt64()),
		IsDelegationClient:   mc.IsDelegationClient.ValueBool(),
		ClientClaims:         convertSetToClientClaims(), // TODO
	}
}

func (mc *MachineClientResource) FromElvidApiClient(client *elvidapiclient.MachineClient) {
	mc.Name = types.StringValue(client.ClientName)
	// mc.Scopes = convertStringSliceToSet(client.Scopes)
	mc.TestUserLoginEnabled = types.BoolValue(client.TestUserLoginEnabled)
	mc.AccessTokenLifeTime = types.Int64Value(int64(client.AccessTokenLifeTime))
	mc.IsDelegationClient = types.BoolValue(client.IsDelegationClient)
	// mc.ClientClaims = convertClientClaimsToSet(client.ClientClaims)
	mc.ClientID = types.StringValue(client.ClientId)
}

func convertSetToStringSlice(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}

// func convertStringSliceToSet(slice []string) types.Set {
// 	var elements []attr.Value
// 	for _, v := range slice {
// 		elements = append(elements, types.StringValue(v))
// 	}
// 	return types.Set{ElementType: types.StringType, Elements: elements}
// }

// func convertSetToClientClaims(set types.Set) []elvidapiclient.ClientClaim {
// 	var result []elvidapiclient.ClientClaim
// 	for _, v := range set.Elements() {
// 		claimMap := v.(types.Object)
// 		claimType := claimMap.Attr("type").(types.String).ValueString()
// 		values := convertSetToStringSlice(claimMap.Attr("values").(types.Set))
// 		result = append(result, elvidapiclient.ClientClaim{Type: claimType, Values: values})
// 	}
// 	return result
// }

// func convertClientClaimsToSet(claims []elvidapiclient.ClientClaim) types.Set {
// 	var elements []attr.Value
// 	for _, claim := range claims {
// 		claimMap := map[string]attr.Value{
// 			"type":   types.StringValue(claim.Type),
// 			"values": convertStringSliceToSet(claim.Values),
// 		}
// 		elements = append(elements, types.Object{AttrTypes: map[string]attr.Type{"type": types.StringType, "values": types.SetType{ElementType: types.StringType}}, Attrs: claimMap})
// 	}
// 	return types.Set{ElementType: types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "values": types.SetType{ElementType: types.StringType}}}, Elements: elements}
// }

func convertSetToClientClaims() []elvidapiclient.ClientClaim {

	return []elvidapiclient.ClientClaim{
		{
			Type:   "client_test_type_1",
			Values: []string{"test_value_1", "test_value_2"},
		},
	}
}
