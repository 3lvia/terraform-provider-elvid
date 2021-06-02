package main

import (
	"fmt"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceApiScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiScopeCreateOrUpdate,
		Read:   resourceApiScopeRead,
		Delete: resourceApiScopeDelete,
		Update: resourceApiScopeCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				providerInput := m.(*ElvidProviderInput)
				d.SetId(d.Id())
				d.Set("resource_taint_version", "1")
				d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")
				err := resourceApiScopeRead(d, m)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the API scope",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     false,
				Description: "A description of what this API scope is used for. Please include information about what it gives access to, and in what way it differs from similar API scopes, if any.",
			},
			"user_claims": {
				Type:     schema.TypeSet,
				Required: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
				// TODO: validere mot godkjente claims? Det holder vel at det feiler ved apply?
			},
			"allow_machine_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the API scope is intended for machine clients (allow_machine_clients and allow_user_clients are mutually exclusive)",
			},
			// TODO: validering av at bare en av disse er satt. Eller skal vi mappe om til enum? Burde også ha validering av at bare en av disse er satt i API'et.
			"allow_user_clients": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the API scope is intended for user clients (allow_machine_clients and allow_user_clients are mutually exclusive)",
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

func resourceApiScopeCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	apiScopeInput := ReadApiScopeFromResourceData(d)
	providerInput := m.(*ElvidProviderInput)
	apiScope, err := elvidapiclient.CreateOrUpdateApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, apiScopeInput)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	d.SetId(apiScope.Name)
	// TODO: hvilke felter skal lagres her? Read-funksjonen setter alt som trengs, utenom Id?
	// Er token_endpoint kun for info?
	// d.Set("client_id", apiScope.ClientId)
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	return resourceApiScopeRead(d, m)
}

func resourceApiScopeRead(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)

	apiScope, err := elvidapiclient.ReadApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
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

func resourceApiScopeDelete(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)

	err := elvidapiclient.DeleteApiScope(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
}

// ClientName:           d.Get("name").(string),
// Scopes:               getStringArrayFromResourceSet(d, "scopes"),
// TestUserLoginEnabled: d.Get("test_user_login_enabled").(bool),
// AccessTokenLifeTime:  d.Get("access_token_life_time").(int),

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
