package main

import (
	"context"
	"fmt"

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
		Schema: map[string]*schema.Schema{
			// The id field has to be Optional and Computed. So we have the name field in addition to the id field even if they'll have the same value.
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
				Description: "User claims that are included in the token when logging in with a machine/user client that has this API scope. (The token will include the superset of claims from all granted scopes).",
				// We don't validate claims in Terraform. It is done by the Elvid API, and apply will fail if invalid scopes are used.
			},
			"allow_machine_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for machine clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
				// We don't validate allow_machine_clients XOR allow_user_clients in Terraform (haven't found a way to access other fields in ValidateFunc). It is validated in the Elvid API.
			},
			"allow_user_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the API scope is intended for user clients (allow_machine_clients and allow_user_clients are mutually exclusive, and one of them has to be true)",
				// We don't validate allow_machine_clients XOR allow_user_clients in Terraform (haven't found a way to access other fields in ValidateFunc). It is validated in the Elvid API.
			},
			"resource_taint_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "1",
				Description: "A change in value for this field will force recreating the resource",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
	d.Set("name", apiScope.Name) // Må sette denne også, for resourceApiScopeRead bruker den

	return resourceApiScopeRead(ctx, d, m)
}

func resourceApiScopeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*ElvidProviderInput)

	var localDiags diag.Diagnostics
	localDiags = append(localDiags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Info: reading api client for resource",
		Detail:   fmt.Sprintf("Debug info: this is the resource converted to ApiScope instance: %+v", ReadApiScopeFromResourceData(d)),
	})

	apiScope, diags := elvidapiclient.ReadApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Get("name").(string))

	if diags != nil {
		diags = append(localDiags, diags...)
		return diags
	}

	if apiScope == nil {
		// If the api scope is not found, let terraform know it does not exist.
		d.SetId("")
		return nil
	}

	d.SetId(apiScope.Name)
	d.Set("name", apiScope.Name)
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
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		UserClaims:          getStringArrayFromResourceSet(d, "user_claims"),
		AllowMachineClients: d.Get("allow_machine_clients").(bool),
		AllowUserClients:    d.Get("allow_user_clients").(bool),
	}
	return apiScope
}
