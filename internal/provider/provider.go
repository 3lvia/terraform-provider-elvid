package provider

import (
	"context"

	"github.com/3lvia/terraform-provider-elvid/internal/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (p *ElvidProviderInput) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elvid"
}

func (p *ElvidProviderInput) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:    true,
				Description: "Azure tenant id",
			},
			"terraform_sp_client_id": schema.StringAttribute{
				Required:    true,
				Description: "The Client ID for terraform service principal",
			},
			"terraform_sp_client_secret": schema.StringAttribute{
				Required:    true,
				Description: "The Client Secret for terraform service principal",
			},
			"environment": schema.StringAttribute{
				Required: true,
			},
			"override_elvid_authority": schema.StringAttribute{
				Optional:    true,
				Description: "elvid_authority is as default set based on environment ( var.environment == 'prod' ? 'https://elvid.elvia.io' : 'https://elvid.test-elvia.io'). Use this to override the default.",
			},
			"run_hashed_secret_validation": schema.BoolAttribute{
				Optional: true,
				// Default:     booldefault.StaticBool(true), must be done in Configure
				Description: "Only turn this off temporarily to recreate a client secret if the hashed_secret_validation fails",
			},
		},
	}
}

func (p *ElvidProviderInput) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config struct {
		TenantID                  types.String `tfsdk:"tenant_id"`
		TerraformSpClientID       types.String `tfsdk:"terraform_sp_client_id"`
		TerraformSpClientSecret   types.String `tfsdk:"terraform_sp_client_secret"`
		Environment               types.String `tfsdk:"environment"`
		OverrideElvidAuthority    types.String `tfsdk:"override_elvid_authority"`
		RunHashedSecretValidation types.Bool   `tfsdk:"run_hashed_secret_validation"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	elvidAuthority := getElvIDAuthority(config.OverrideElvidAuthority.ValueString(), config.Environment.ValueString())

	accessTokenAD, err := elvidapiclient.GetAccessTokenAD(config.TenantID.ValueString(), config.TerraformSpClientID.ValueString(), config.TerraformSpClientSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to get access token", err.Error())
		return
	}

	runHashedSecretValidation := config.RunHashedSecretValidation.ValueBool()
	if config.RunHashedSecretValidation.IsNull() {
		runHashedSecretValidation = true
	}

	providerInput := &ElvidProviderInput{
		TenantId:                  config.TenantID.ValueString(),
		AccessTokenAD:             accessTokenAD,
		ElvIDAuthority:            elvidAuthority,
		RunHashedSecretValidation: runHashedSecretValidation,
	}

	resp.DataSourceData = providerInput
	resp.ResourceData = providerInput
}

func getElvIDAuthority(overrideElvidAuthority string, environment string) string {
	if overrideElvidAuthority != "" {
		return overrideElvidAuthority
	} else if environment == "prod" {
		return "https://elvid.elvia.io"
	} else {
		return "https://elvid.test-elvia.io"
	}
}

type ElvidProviderInput struct {
	TenantId                  string
	AccessTokenAD             string
	ElvIDAuthority            string
	RunHashedSecretValidation bool
}

// DataSources implements provider.Provider.
func (p *ElvidProviderInput) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources implements provider.Provider.
func (p *ElvidProviderInput) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMachineClientResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ElvidProviderInput{}
	}
}
