package main

import (
	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Azure tenant id",
				ValidateFunc: validation.IsUUID,
			},
			"terraform_sp_client_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Client ID for terraform service principal",
				ValidateFunc: validation.IsUUID,
			},
			"terraform_sp_client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Client Secret for terraform service principal",
			},
			"environment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"override_elvid_authority": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "elvid_authority is as default set based on environment ( var.environment == 'prod' ? 'https://elvid.elvia.io' : 'https://elvid.test-elvia.io'). Use this to override the default.",
			},
			"run_hashed_secret_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Only turn this off temperary to recreate a client secret if the hashed_secret_validation fails",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"elvid_clientsecret":  resourceClientSecret(),
			"elvid_machineclient": resourceMachineClient(),
			"elvid_userclient":    resourceUserClient(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func getElvIDAuthoriry(overrideElvidAuthority string, environemnt string) string {
	if overrideElvidAuthority != "" {
		return overrideElvidAuthority
	} else if environemnt == "prod" {
		return "https://elvid.elvia.io"
	} else {
		return "https://elvid.test-elvia.io"
	}
}

// note that the interface{} response from this  will be accessable from the resources as m interface{}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	tenantID := d.Get("tenant_id").(string)
	environment := d.Get("environment").(string)
	terraformSpClientId := d.Get("terraform_sp_client_id").(string)
	terraformSpClientSecret := d.Get("terraform_sp_client_secret").(string)
	overrideElvidAuthority := d.Get("override_elvid_authority").(string)

	elvidAuthority := getElvIDAuthoriry(overrideElvidAuthority, environment)

	runHashedSecretValidation := d.Get("run_hashed_secret_validation").(bool)
	accessTokenAD, err := elvidapiclient.GetAccessTokenAD(tenantID, terraformSpClientId, terraformSpClientSecret)

	if err != nil {
		return nil, err
	}
	providerInput := &ElvidProviderInput{tenantID, accessTokenAD, elvidAuthority, runHashedSecretValidation}
	return providerInput, nil
}

func getStringArrayFromResourceSet(d *schema.ResourceData, name string) []string {
	rawList := d.Get(name).(*schema.Set).List()
	stringList := make([]string, len(rawList))
	for i, v := range rawList {
		stringList[i] = v.(string)
	}
	return stringList
}

type ElvidProviderInput struct {
	TenantId                  string
	AccessTokenAD             string
	ElvIDAuthority            string
	RunHashedSecretValidation bool
}
