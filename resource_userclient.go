package main

import (
	"fmt"
	"strconv"

	"github.com/3lvia/terraform-provider-elvid/elvidapiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserClientCreate,
		Read:   resourceUserClientRead,
		Delete: resourceUserClientDelete,
		Update: resourceUserClientUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				d.SetId(d.Id())
				d.Set("resource_taint_version", "1")
				err := resourceUserClientRead(d, m)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"client_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The full name of the client. Note that ClientNameLanguageKey can be used to get a seperate language dependent name.",
			},
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The scopes the client can request, note that scopes are auto approved in test and must be approved by an admin in production.",
			},
			"domains": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The domains the client is found in. Used for Cors AllowedCorsOrigins, RedirectUris and PostLogoutRedirectUris.",
			},
			"redirect_uri_paths": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The path part of a RedirectUri, each of these will be combined with each of the domains. ElvID is only allowed to send the user back to the client with one of these uris.",
			},
			"post_logout_redirect_uri_paths": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The path part of a PostLogoutRedirectUri, each of these will be combined with each of the domains. After logout ElvID is only allowed to send the user back to the client with one of these uris.",
			},
			"bankid_login_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to allow user to log in with BankID.",
			},
			"local_login_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to allow user to log in with local account e.g. email and password. Note this is not for the work-related ad-email.",
			},
			"elvia_ad_login_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to allow user to log in with Elvia AD.",
			},
			"test_user_login_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When this is enabled it's possible to login with a test user to this client, by using local login when authenticating the user in the gui. See ElvID space on confluence for details.",
			},
			"require_client_secret": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "ClientSecret are typically not required in userclients, as most of them run on the users hardware and cannot keep a secret. A ClientSecret can be used safely from the backend of a dynamic webpage like ASP.NET Core MVC.",
			},
			"access_token_life_time": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Number of seconds before an access token expires. Default and max is 3600 seconds (1 hour)",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 3600 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 3600 inclusive, got: %d", key, v))
					}
					return
				},
			},
			"always_include_user_claims_in_id_token": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When requesting both an id token and access token, should the user claims always be added to the id token instead of requiring the client to use the userinfo endpoint?",
			},
			"client_name_language_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Use this to get a language dependent seperate name for the client. That name could be used instead of the client name in places like the BackToClient-button that is showed for an signed in user in elvid. Note that a corresponding key/value for 'ClientName{client_name_language_key}' must also exist in elvid's language files. Eg language.nb.json --> key: 'ClientNameMinSide', value: 'MinSide'",
			},
			"allow_use_of_refresh_tokens": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "This will enable the use of refreshtokens, and the scope offline_access. Talk with GlueTeam before using this.",
			},
			"one_time_usage_for_refresh_tokens": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "This have no use if allow_use_of_refresh_tokens = false. When using a OneTime RefreshToken the token endpoint response includes a new RefreshToken that should be used for the next token request. Set this to false to get a reusable RefreshToken. OneTime RefreshTokens are safer as an attacker cannot replay a request to the token endoint with the used RefreshToken or othervise use an already used RefreshToken. For the field RefreshTokenUsage in the elvid DB 0 means ReUse and 1 means SingleUse.",
			},
			"refresh_token_life_time": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2592000,
				Description: "This have no use if allow_use_of_refresh_tokens = false. Absolute number of seconds before a refresh token expires. Note that a refresh token can also be revoked. Default is 2592000 seconds (30 days), max is 31556926 seconds (1 year).",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 || v > 31556926 {
						errs = append(errs, fmt.Errorf("%q must be between 1 and 31556926 inclusive, got: %d", key, v))
					}
					return
				},
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
		},
	}
}

func resourceUserClientCreate(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)
	userClient := ReadUserClientFromResourceData(d)
	createdUserClient, err := elvidapiclient.CreateUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, userClient)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	d.SetId(strconv.Itoa(createdUserClient.Id))
	d.Set("client_id", createdUserClient.ClientId)
	return resourceUserClientRead(d, m)
}

func resourceUserClientRead(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)

	userClient, err := elvidapiclient.ReadUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	if userClient == nil {
		// If the client is not found, let terraform know it does not exist.
		d.SetId("")
		return nil
	}

	d.Set("client_id", userClient.ClientId)
	d.Set("scopes", userClient.Scopes)
	d.Set("domains", userClient.Domains)
	d.Set("redirect_uri_paths", userClient.RedirectUriPaths)
	d.Set("post_logout_redirect_uri_paths", userClient.PostLogoutRedirectUriPaths)
	d.Set("bankid_login_enabled", userClient.BankIDLoginEnabled)
	d.Set("local_login_enabled", userClient.LocalLoginEnabled)
	d.Set("elvia_ad_login_enabled", userClient.ElviaADLoginEnabled)
	d.Set("test_user_login_enabled", userClient.TestUserLoginEnabled)
	d.Set("require_client_secret", userClient.RequireClientSecret)
	d.Set("access_token_life_time", userClient.AccessTokenLifetime)
	d.Set("always_include_user_claims_in_id_token", userClient.AlwaysIncludeUserClaimsInIdToken)
	d.Set("client_name_language_key", userClient.ClientNameLanguageKey)
	d.Set("client_name", userClient.ClientName)
	d.Set("allow_use_of_refresh_tokens", userClient.AllowUseOfRefreshTokens)
	d.Set("one_time_usage_for_refresh_tokens", userClient.OneTimeUsageForRefreshTokens)
	d.Set("refresh_token_life_time", userClient.RefreshTokensLifeTime)

	return nil
}

func resourceUserClientUpdate(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)
	userClient := ReadUserClientFromResourceData(d)
	userClient.Id, _ = strconv.Atoi(d.Id())
	userClient.ClientId = d.Get("client_id").(string)

	err := elvidapiclient.UpdateUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, userClient)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	return resourceUserClientRead(d, m)
}

func resourceUserClientDelete(d *schema.ResourceData, m interface{}) error {
	providerInput := m.(*ElvidProviderInput)
	err := elvidapiclient.DeleteUserClient(providerInput.ElvIDAuthority, providerInput.AccessTokenAD, d.Id())

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}
	return nil
}

func ReadUserClientFromResourceData(d *schema.ResourceData) *elvidapiclient.UserClient {
	userClient := &elvidapiclient.UserClient{
		ClientName:                       d.Get("client_name").(string),
		Scopes:                           getStringArrayFromResourceSet(d, "scopes"),
		Domains:                          getStringArrayFromResourceSet(d, "domains"),
		RedirectUriPaths:                 getStringArrayFromResourceSet(d, "redirect_uri_paths"),
		PostLogoutRedirectUriPaths:       getStringArrayFromResourceSet(d, "post_logout_redirect_uri_paths"),
		BankIDLoginEnabled:               d.Get("bankid_login_enabled").(bool),
		LocalLoginEnabled:                d.Get("local_login_enabled").(bool),
		ElviaADLoginEnabled:              d.Get("elvia_ad_login_enabled").(bool),
		TestUserLoginEnabled:             d.Get("test_user_login_enabled").(bool),
		RequireClientSecret:              d.Get("require_client_secret").(bool),
		AccessTokenLifetime:              d.Get("access_token_life_time").(int),
		AlwaysIncludeUserClaimsInIdToken: d.Get("always_include_user_claims_in_id_token").(bool),
		ClientNameLanguageKey:            d.Get("client_name_language_key").(string),
		AllowUseOfRefreshTokens:          d.Get("allow_use_of_refresh_tokens").(bool),
		OneTimeUsageForRefreshTokens:     d.Get("one_time_usage_for_refresh_tokens").(bool),
		RefreshTokensLifeTime:            d.Get("refresh_token_life_time").(int),
	}
	return userClient
}
