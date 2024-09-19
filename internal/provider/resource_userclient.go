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

func NewUserClientResource() resource.Resource {
	return &UserClientResource{}
}

func (r *UserClientResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "elvid_userclient"
}

func (r *UserClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		providerInput = req.ProviderData.(*ElvidProviderInput) // move to provider.go
	}
}

func (r *UserClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_name": schema.StringAttribute{
				Required:    true,
				Description: "The full name of the client. Note that ClientNameLanguageKey can be used to get a separate language dependent name.",
			},
			"scopes": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
			},
			"domains": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The domains the client is found in. Used for Cors AllowedCorsOrigins, RedirectUris and PostLogoutRedirectUris.",
			},
			"redirect_uri_paths": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The path part of a RedirectUri, each of these will be combined with each of the domains. ElvID is only allowed to send the user back to the client with one of these uris.",
			},
			"post_logout_redirect_uri_paths": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "The path part of a PostLogoutRedirectUri, each of these will be combined with each of the domains. After logout ElvID is only allowed to send the user back to the client with one of these uris.",
			},
			"idporten_login_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable to allow user to log in with ID-porten.",
			},
			"local_login_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable to allow user to log in with local account e.g. email and password. Note this is not for the work-related ad-email.",
			},
			"elvia_ad_login_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable to allow user to log in with Elvia AD.",
			},
			"test_user_login_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When this is enabled it's possible to login with a test user to this client, by using local login when authenticating the user in the gui. See ElvID space on confluence for details.",
			},
			"require_client_secret": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "ClientSecret are typically not required in userclients, as most of them run on the users hardware and cannot keep a secret. A ClientSecret can be used safely from the backend of a dynamic webpage like ASP.NET Core MVC.",
			},
			"access_token_life_time": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3600),
				Description: "Number of seconds before an access token expires. Default and max is 3600 seconds (1 hour)",
			},
			"always_include_user_claims_in_id_token": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "When requesting both an id token and access token, should the user claims always be added to the id token instead of requiring the client to use the userinfo endpoint?",
			},
			"client_name_language_key": schema.StringAttribute{
				Optional:    true,
				Description: "Use this to get a language dependent separate name for the client. That name could be used instead of the client name in places like the BackToClient-button that is showed for an signed in user in elvid. Note that a corresponding key/value for 'ClientName{client_name_language_key}' must also exist in elvid's language files. Eg language.nb.json --> key: 'ClientNameMinSide', value: 'MinSide'",
			},
			"allow_use_of_refresh_tokens": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "This will enable the use of refreshtokens, and the scope offline_access. Talk with GlueTeam before using this.",
			},
			"one_time_usage_for_refresh_tokens": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "This have no use if allow_use_of_refresh_tokens = false. When using a OneTime RefreshToken the token endpoint response includes a new RefreshToken that should be used for the next token request. Set this to false to get a reusable RefreshToken. For the field RefreshTokenUsage in the elvid DB 0 means ReUse and 1 means SingleUse.",
			},
			"refresh_token_life_time": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(2592000),
				Description: "This have no use if allow_use_of_refresh_tokens = false. Absolute number of seconds before a refresh token expires. Note that a refresh token can also be revoked. Default is 2592000 seconds (30 days), max is 31556926 seconds (1 year).",
			},
			"client_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The cliend_id of the client, used during client_credentials auth. Note this is different from the (entity) id of the client",
			},
			"resource_taint_version": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Computed:    true,
				Default:     stringdefault.StaticString("1"),
				Description: "A change in value for this field will force recreating the resource",
			},
			// "client_properties": schema.SetNestedAttribute{
			// 	Optional:    true,
			// 	Description: "Used this to set other key-value(s) properties on a client. The allowed keys to set here must be whitelisted in elvid.",
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"type": schema.StringAttribute{
			// 				Required: true,
			// 			},
			// 			"values": schema.SetAttribute{
			// 				ElementType: types.StringType,
			// 				Required:    true,
			// 			},
			// 		},
			// 	},
			// },
		},
	}
}

func (r *UserClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserClientResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userClientInput := plan.ToElvidApiClientInput()
	createdUserClient, err := elvidapiclient.CreateUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, userClientInput)
	if err != nil {
		resp.Diagnostics.AddError("Creating userclient resulted in an error", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdUserClient.Id))
	plan.ClientID = types.StringValue(createdUserClient.ClientId)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *UserClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserClientResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userClient, err := elvidapiclient.ReadUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Reading userclient resulted in an error", err.Error())
		return
	}

	if userClient == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.FromElvidApiClient(userClient)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *UserClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserClientResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	userClientInput := plan.ToElvidApiClientInput()
	userClientInput.Id, _ = strconv.Atoi(plan.ID.ValueString())

	err := elvidapiclient.UpdateUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, userClientInput)
	if err != nil {
		resp.Diagnostics.AddError("Updating userclient resulted in an error", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *UserClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserClientResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := elvidapiclient.DeleteUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Deleting userclient resulted in an error", err.Error())
	}
}

type UserClientResource struct {
	ID                               types.String `tfsdk:"id"`
	ClientName                       types.String `tfsdk:"client_name"`
	Scopes                           types.Set    `tfsdk:"scopes"`
	Domains                          types.Set    `tfsdk:"domains"`
	RedirectUriPaths                 types.Set    `tfsdk:"redirect_uri_paths"`
	PostLogoutRedirectUriPaths       types.Set    `tfsdk:"post_logout_redirect_uri_paths"`
	IdPortenLoginEnabled             types.Bool   `tfsdk:"idporten_login_enabled"`
	LocalLoginEnabled                types.Bool   `tfsdk:"local_login_enabled"`
	ElviaADLoginEnabled              types.Bool   `tfsdk:"elvia_ad_login_enabled"`
	TestUserLoginEnabled             types.Bool   `tfsdk:"test_user_login_enabled"`
	RequireClientSecret              types.Bool   `tfsdk:"require_client_secret"`
	AccessTokenLifetime              types.Int64  `tfsdk:"access_token_life_time"`
	AlwaysIncludeUserClaimsInIdToken types.Bool   `tfsdk:"always_include_user_claims_in_id_token"`
	ClientNameLanguageKey            types.String `tfsdk:"client_name_language_key"`
	AllowUseOfRefreshTokens          types.Bool   `tfsdk:"allow_use_of_refresh_tokens"`
	OneTimeUsageForRefreshTokens     types.Bool   `tfsdk:"one_time_usage_for_refresh_tokens"`
	RefreshTokensLifeTime            types.Int64  `tfsdk:"refresh_token_life_time"`
	ClientID                         types.String `tfsdk:"client_id"`
	ResourceTaintVersion             types.String `tfsdk:"resource_taint_version"`
	// ClientProperties                 types.Set    `tfsdk:"client_properties"`
}

func (uc *UserClientResource) ToElvidApiClientInput() *elvidapiclient.UserClient {
	return &elvidapiclient.UserClient{
		ClientName:                       uc.ClientName.ValueString(),
		Scopes:                           convertSetToStringArray(uc.Scopes),
		Domains:                          convertSetToStringArray(uc.Domains),
		RedirectUriPaths:                 convertSetToStringArray(uc.RedirectUriPaths),
		PostLogoutRedirectUriPaths:       convertSetToStringArray(uc.PostLogoutRedirectUriPaths),
		IdPortenLoginEnabled:             uc.IdPortenLoginEnabled.ValueBool(),
		LocalLoginEnabled:                uc.LocalLoginEnabled.ValueBool(),
		ElviaADLoginEnabled:              uc.ElviaADLoginEnabled.ValueBool(),
		TestUserLoginEnabled:             uc.TestUserLoginEnabled.ValueBool(),
		RequireClientSecret:              uc.RequireClientSecret.ValueBool(),
		AccessTokenLifetime:              int(uc.AccessTokenLifetime.ValueInt64()),
		AlwaysIncludeUserClaimsInIdToken: uc.AlwaysIncludeUserClaimsInIdToken.ValueBool(),
		ClientNameLanguageKey:            uc.ClientNameLanguageKey.ValueString(),
		AllowUseOfRefreshTokens:          uc.AllowUseOfRefreshTokens.ValueBool(),
		OneTimeUsageForRefreshTokens:     uc.OneTimeUsageForRefreshTokens.ValueBool(),
		RefreshTokensLifeTime:            int(uc.RefreshTokensLifeTime.ValueInt64()),
		// ClientProperties:                 convertSetToClientProperties(uc.ClientProperties),
	}
}

func (uc *UserClientResource) FromElvidApiClient(client *elvidapiclient.UserClient) {
	uc.ClientName = types.StringValue(client.ClientName)
	// uc.Scopes = convertStringArrayToSet(client.Scopes)
	// uc.Domains = convertStringArrayToSet(client.Domains)
	// uc.RedirectUriPaths = convertStringArrayToSet(client.RedirectUriPaths)
	// uc.PostLogoutRedirectUriPaths = convertStringArrayToSet(client.PostLogoutRedirectUriPaths)
	uc.IdPortenLoginEnabled = types.BoolValue(client.IdPortenLoginEnabled)
	uc.LocalLoginEnabled = types.BoolValue(client.LocalLoginEnabled)
	uc.ElviaADLoginEnabled = types.BoolValue(client.ElviaADLoginEnabled)
	uc.TestUserLoginEnabled = types.BoolValue(client.TestUserLoginEnabled)
	uc.RequireClientSecret = types.BoolValue(client.RequireClientSecret)
	uc.AccessTokenLifetime = types.Int64Value(int64(client.AccessTokenLifetime))
	uc.AlwaysIncludeUserClaimsInIdToken = types.BoolValue(client.AlwaysIncludeUserClaimsInIdToken)
	uc.ClientNameLanguageKey = types.StringValue(client.ClientNameLanguageKey)
	uc.AllowUseOfRefreshTokens = types.BoolValue(client.AllowUseOfRefreshTokens)
	uc.OneTimeUsageForRefreshTokens = types.BoolValue(client.OneTimeUsageForRefreshTokens)
	uc.RefreshTokensLifeTime = types.Int64Value(int64(client.RefreshTokensLifeTime))
	uc.ClientID = types.StringValue(client.ClientId)
	// uc.ClientProperties = convertClientPropertiesToSet(client.ClientProperties)
}

func convertSetToStringArray(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		result = append(result, v.(types.String).ValueString())
	}
	return result
}
