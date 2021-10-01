package main

import (
	"fmt"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMachineClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineClientCreate,
		Read:   resourceMachineClientRead,
		Delete: resourceMachineClientDelete,
		Update: resourceMachineClientUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				providerInput := m.(*ElvidProviderInput)
				d.SetId(d.Id())
				d.Set("resource_taint_version", "1")
				d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")
				err := resourceMachineClientRead(d, m)
				return []*schema.ResourceData{d}, err
			},
		},
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
		},
	}
}

func resourceMachineClientCreate(d *schema.ResourceData, m interface{}) error {
	machineClientInput := ReadMachineClientFromResourceData(d)
	providerInput := m.(*ElvidProviderInput)
	machineClient, err := elvidapiclient.CreateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	d.SetId(strconv.Itoa(machineClient.Id))
	d.Set("client_id", machineClient.ClientId)
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	return resourceMachineClientRead(d, m)
}

func resourceMachineClientRead(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)

	machineClient, err := elvidapiclient.ReadMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	if machineClient == nil {
		// If the client is not found, let terraform know it does not exist.
		d.SetId("")
		return nil
	}

	d.Set("client_id", machineClient.ClientId)
	d.Set("name", machineClient.ClientName)
	d.Set("access_token_life_time", machineClient.AccessTokenLifeTime)
	d.Set("test_user_login_enabled", machineClient.TestUserLoginEnabled)
	d.Set("is_delegation_client", machineClient.IsDelegationClient)

	d.Set("scopes", machineClient.Scopes)

	return nil
}

func resourceMachineClientUpdate(d *schema.ResourceData, m interface{}) error {
	machineClientInput := ReadMachineClientFromResourceData(d)
	machineClientInput.Id, _ = strconv.Atoi(d.Id())
	providerInput := m.(*ElvidProviderInput)

	_, err := elvidapiclient.UpdateMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, machineClientInput)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	d.Set("token_endpoint", providerInput.ElvIDAuthority+"/connect/token")

	return resourceMachineClientRead(d, m)
}

func resourceMachineClientDelete(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)

	err := elvidapiclient.DeleteMachineClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
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
