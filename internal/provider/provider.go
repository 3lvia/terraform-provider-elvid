package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

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
				Optional:    true,
				Sensitive:   true,
				Description: "The Client Secret for the terraform service principal. Leave empty to log in with a client assertion instead (see terraform_sp_client_assertion).",
			},
			"terraform_sp_client_assertion": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "An OIDC token issued to the run by its platform (Scalr, Terraform Cloud, GitHub Actions) that the service principal's app registration trusts as a federated identity credential. Exchanged for an Entra access token with client_assertion, so no client secret exists. Defaults to ARM_OIDC_TOKEN, or the contents of the file named by ARM_OIDC_TOKEN_FILE_PATH, the same variables the azurerm provider reads.",
			},
			"environment": schema.StringAttribute{
				Required: true,
			},
			"override_elvid_authority": schema.StringAttribute{
				Optional:    true,
				Description: "elvid_authority is as default set based on environment ( var.environment == 'prod' ? 'https://elvid.elvia.io' : 'https://elvid.test-elvia.io'). Use this to override the default.",
			},
			"run_hashed_secret_validation": schema.BoolAttribute{
				Optional:    true,
				Description: "Only turn this off temporarily to recreate a client secret if the hashed_secret_validation fails",
			},
		},
	}
}

func (p *ElvidProviderInput) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config struct {
		TenantID                   types.String `tfsdk:"tenant_id"`
		TerraformSpClientID        types.String `tfsdk:"terraform_sp_client_id"`
		TerraformSpClientSecret    types.String `tfsdk:"terraform_sp_client_secret"`
		TerraformSpClientAssertion types.String `tfsdk:"terraform_sp_client_assertion"`
		Environment                types.String `tfsdk:"environment"`
		OverrideElvidAuthority     types.String `tfsdk:"override_elvid_authority"`
		RunHashedSecretValidation  types.Bool   `tfsdk:"run_hashed_secret_validation"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	elvidAuthority := getElvIDAuthority(config.OverrideElvidAuthority.ValueString(), config.Environment.ValueString())

	clientAssertion, err := resolveClientAssertion(config.TerraformSpClientAssertion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read client assertion", err.Error())
		return
	}
	if clientAssertion == "" && config.TerraformSpClientSecret.ValueString() == "" {
		resp.Diagnostics.AddError("No credentials for the terraform service principal",
			"Set terraform_sp_client_assertion (or ARM_OIDC_TOKEN / ARM_OIDC_TOKEN_FILE_PATH) to log in with a federated token, or terraform_sp_client_secret to log in with a secret.")
		return
	}

	accessTokenAD, err := elvidapiclient.GetAccessTokenAD(config.TenantID.ValueString(), config.TerraformSpClientID.ValueString(), config.TerraformSpClientSecret.ValueString(), clientAssertion)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get access token", err.Error())
		return
	}

	runHashedSecretValidation := config.RunHashedSecretValidation.ValueBool()
	if config.RunHashedSecretValidation.IsNull() {
		runHashedSecretValidation = true
	}

	elvidProviderInput := &ElvidProviderInput{
		TenantId:                  config.TenantID.ValueString(),
		AccessTokenAD:             accessTokenAD,
		ElvIDAuthority:            elvidAuthority,
		RunHashedSecretValidation: runHashedSecretValidation,
	}

	resp.DataSourceData = elvidProviderInput
	resp.ResourceData = elvidProviderInput
}

// resolveClientAssertion returns the explicitly configured assertion, else ARM_OIDC_TOKEN, else the contents of
// the file in ARM_OIDC_TOKEN_FILE_PATH; empty when none of them is set.
func resolveClientAssertion(configured string) (string, error) {
	if configured != "" {
		return configured, nil
	}
	if token := os.Getenv("ARM_OIDC_TOKEN"); token != "" {
		return token, nil
	}
	if path := os.Getenv("ARM_OIDC_TOKEN_FILE_PATH"); path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ARM_OIDC_TOKEN_FILE_PATH: %w", err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return "", nil
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
		MachineClientResourceSetup,
		ClientSecretResourceSetup,
		UserClientResourceSetup,
		ApiScopeResourceSetup,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ElvidProviderInput{}
	}
}
