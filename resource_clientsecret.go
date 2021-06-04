package main

import (
	"fmt"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClientSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceClientSecretCreate,
		Read:   resourceClientSecretRead,
		Delete: resourceClientSecretDelete,

		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The (entity) id of the client. Note that this should be set to client.id and not client.client_id",
			},
			"resource_taint_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "1",
				Description: "A change in value for this field will force recreating the resource",
			},
			"secret_value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"hashed_value_starts_with": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceClientSecretCreate(d *schema.ResourceData, m interface{}) error {
	clientId := d.Get("client_id").(string)
	providerInput := m.(*ElvidProviderInput)

	clientSecret, err := elvidapiclient.CreateClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, clientId)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	d.SetId(strconv.Itoa(clientSecret.Id))
	d.Set("secret_value", clientSecret.Value)
	d.Set("hashed_value_starts_with", clientSecret.HashedValueStartsWith)

	return resourceClientSecretRead(d, m)
}

func resourceClientSecretRead(d *schema.ResourceData, m interface{}) error {
	clientId := d.Get("client_id").(string)
	clientSecretId := d.Id()

	providerInput := m.(*ElvidProviderInput)

	clientSecret, err := elvidapiclient.ReadClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, clientId, clientSecretId)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	if clientSecret == nil {
		// If the clientSecret is not found, let terraform know it does not exist.
		d.SetId("")
		return nil
	}

	if providerInput.RunHashedSecretValidation && clientSecret.HashedValueStartsWith != d.Get("hashed_value_starts_with").(string) {
		return fmt.Errorf("The secret was found but hashed_value_starts_with was not as expected, meaning the secret must have changed without terraform knowning. Either fix the state or secret manuelly, or recreate the secret. Recreate the secret by TEMPERARY setting run_hashed_secret_validation = false in the provider and changing resource_taint_version on the resource")
	}

	return nil
}

func resourceClientSecretDelete(d *schema.ResourceData, m interface{}) error {
	clientId := d.Get("client_id").(string)
	clientSecretId := d.Id()

	providerInput := m.(*ElvidProviderInput)

	err := elvidapiclient.DeleteClientSecret(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, clientId, clientSecretId)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
}
