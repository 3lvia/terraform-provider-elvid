package main

import (
	"context"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachineClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMachineClientCreate,
		ReadContext:   resourceMachineClientRead,
		DeleteContext: resourceMachineClientDelete,
		UpdateContext: resourceMachineClientUpdate,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the client",
			},
			"test_user_login_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When this is enabled it's possible to mechanically login with a test user to this client. This is done by using grant-type password for the token endpoint. See ElvID space on confluence for details",
			},
			"is_delegation_client": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When this is enabled the client can only use the delegation grant type. This is used when a already logged inn user will create a long-lived delegation access_token",
			},
			"access_token_life_time": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Number of seconds before an access token expires. Default and maximum is 3600 seconds (1 hour) for regular machine clients.",
			},
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cliend_id of the client, used during client_credentials auth. Note this is different from the (entity) id of the client",
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
			"client_claims": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"claims": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceMachineClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	machineClientInput := ReadMachineClientFromResourceData(d)
	providerInput := m.(*ElvidProviderInput)
	machineClient, err := elvidapiclient.CreateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)
	var diags diag.Diagnostics

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Creating machineclient resulted in an error",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(strconv.Itoa(machineClient.Id))
	d.Set("client_id", machineClient.ClientId)
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	resourceMachineClientRead(ctx, d, m)
	return diags
}

func resourceMachineClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*ElvidProviderInput)

	machineClient, err := elvidapiclient.ReadMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	var diags diag.Diagnostics

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Reading machineclient resulted in an error",
			Detail:   err.Error(),
		})
		return diags
	}

	if machineClient == nil {
		// If the client is not found, let terraform know it does not exist.
		d.SetId("")
		return diags
	}

	d.Set("client_id", machineClient.ClientId)
	d.Set("name", machineClient.ClientName)
	d.Set("access_token_life_time", machineClient.AccessTokenLifeTime)
	d.Set("test_user_login_enabled", machineClient.TestUserLoginEnabled)
	d.Set("is_delegation_client", machineClient.IsDelegationClient)

	d.Set("scopes", machineClient.Scopes)

	return diags
}

func resourceMachineClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	machineClientInput := ReadMachineClientFromResourceData(d)
	machineClientInput.Id, _ = strconv.Atoi(d.Id())
	providerInput := m.(*ElvidProviderInput)

	_, err := elvidapiclient.UpdateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "WarningS",
		Detail:   "WarningD",
	})

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Updating machineclient resulted in an error",
			Detail:   err.Error(),
		})
		return diags
	}
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	resourceMachineClientRead(ctx, d, m)
	return diags
}

func resourceMachineClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	providerInput := m.(*ElvidProviderInput)
	var diags diag.Diagnostics

	err := elvidapiclient.DeleteMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Deleting machineclient resulted in an error",
			Detail:   err.Error(),
		})
	}
	return diags
}

func ReadMachineClientFromResourceData(d *schema.ResourceData) *elvidapiclient.MachineClient {
	machineClient := &elvidapiclient.MachineClient{
		ClientName:           d.Get("name").(string),
		Scopes:               getStringArrayFromResourceSet(d, "scopes"),
		TestUserLoginEnabled: d.Get("test_user_login_enabled").(bool),
		AccessTokenLifeTime:  d.Get("access_token_life_time").(int),
		IsDelegationClient:   d.Get("is_delegation_client").(bool),
	}
	return machineClient
}
