package main

import (
	"context"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiScope() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiScopeCreateOrUpdate,
		ReadContext:   resourceApiScopeRead,
		DeleteContext: resourceApiScopeDelete,
		UpdateContext: resourceApiScopeCreateOrUpdate,
		// TODO: få denne kompatibel med v2:
		// Importer: &schema.ResourceImporter{
		// 	State: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		// 		providerInput := m.(*ElvidProviderInput)
		// 		d.SetId(d.Id())
		// 		d.Set("resource_taint_version", "1")
		// 		d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")
		// 		diags := resourceApiScopeRead(d, m)

		// 		return []*schema.ResourceData{d}, diags
		// 	},
		// },
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the API scope",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of what this API scope is used for. Please include information about what it gives access to, and in what way it differs from similar API scopes, if any.",
			},
			"user_claims": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
				// TODO: validere mot godkjente claims? Det holder vel at det feiler ved apply?
			},
			"allow_machine_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for machine clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			// TODO: validering av at bare en av disse er satt. Eller skal vi mappe om til enum? Burde også ha validering av at bare en av disse er satt i API'et.
			"allow_user_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for user clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
			},
			"resource_taint_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "1",
				Description: "A change in value for this field will force recreating the resource",
			},
			"token_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApiScopeCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiScopeInput := ReadApiScopeFromResourceData(d)
	providerInput := m.(*ElvidProviderInput)
	apiScope, diags := elvidapiclient.CreateOrUpdateApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, apiScopeInput)

	if diags != nil {
		return diags
	}

	d.SetId(apiScope.Name)
	// TODO: hvilke felter skal lagres her? Read-funksjonen setter alt som trengs, utenom Id?
	// Er token_endpoint kun for info?
	// d.Set("client_id", apiScope.ClientId)
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	return resourceApiScopeRead(ctx, d, m)
}

func resourceApiScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*ElvidProviderInput)

	apiScope, diags := elvidapiclient.ReadApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if diags != nil {
		return diags
	}

	if apiScope == nil {
		// If the api scope is not found, let terraform know it does not exist.
		d.SetId("")
		// TODO: trenger vi ikke nullstille øvrige felter i ResourceData d?
		return nil
	}

	d.Set("name", apiScope.Name) // TODO: trenger vi name i tillegg til Id i Terraform?
	d.Set("description", apiScope.Description)
	d.Set("user_claims", apiScope.UserClaims)
	d.Set("allow_machine_clients", apiScope.AllowMachineClients)
	d.Set("allow_user_clients", apiScope.AllowUserClients)

	return nil
}

func resourceApiScopeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*ElvidProviderInput)

	diags := elvidapiclient.DeleteApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	return diags
}

func ReadApiScopeFromResourceData(d *schema.ResourceData) *elvidapiclient.ApiScope {
	apiScope := &elvidapiclient.ApiScope{
		Name:                d.Id(),
		Description:         d.Get("description").(string),
		UserClaims:          getStringArrayFromResourceSet(d, "user_claims"),
		AllowMachineClients: d.Get("allow_machine_clients").(bool),
		AllowUserClients:    d.Get("allow_user_clients").(bool),
	}
	return apiScope
}
