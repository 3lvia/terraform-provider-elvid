package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

const elvidApiScope = "api://32add519-6233-4728-a7d9-203dd3968436/.default"

// GetAccessTokenAD exchanges the service principal's credentials for an Entra access token for the ElvID API.
// With a client assertion (an OIDC token issued to the run by its platform, e.g. Scalr, and trusted by the app
// registration as a federated identity credential) no password exists anywhere; the client secret is the
// legacy path and is only used when no assertion is given.
func GetAccessTokenAD(tenantId string, clientId string, clientSecret string, clientAssertion string) (string, error) {
	tokenEndpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantId)

	response, err := myClient.PostForm(tokenEndpoint, buildTokenRequest(clientId, clientSecret, clientAssertion))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var adResponse ADResponse
	if err = json.Unmarshal(data, &adResponse); err != nil {
		return "", err
	}
	if adResponse.AccessToken == "" {
		if adResponse.Error != "" {
			return "", fmt.Errorf("Entra returned %s: %s", adResponse.Error, adResponse.ErrorDescription)
		}
		return "", fmt.Errorf("AccessToken not found in response (HTTP %d)", response.StatusCode)
	}
	return adResponse.AccessToken, nil
}

// buildTokenRequest builds the client_credentials form. A client assertion takes precedence over a secret.
func buildTokenRequest(clientId string, clientSecret string, clientAssertion string) url.Values {
	form := url.Values{
		"grant_type": {"client_credentials"},
		"client_id":  {clientId},
		"scope":      {elvidApiScope},
	}
	if clientAssertion != "" {
		form.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
		form.Set("client_assertion", clientAssertion)
	} else {
		form.Set("client_secret", clientSecret)
	}
	return form
}

type ADResponse struct {
	AccessToken      string `json:"access_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
