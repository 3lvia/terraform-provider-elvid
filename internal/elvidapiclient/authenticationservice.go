package elvidapiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

func GetAccessTokenAD(tenantId string, clientId string, clientSecret string) (string, error) {
	tokenEndpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantId)

	response, err := myClient.PostForm(tokenEndpoint, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"scope":         {"api://32add519-6233-4728-a7d9-203dd3968436/.default"}})

	if err != nil {
		return "", err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		var adResponse ADResponse
		err = json.Unmarshal(data, &adResponse)
		if err != nil {
			return "", err
		}

		if adResponse.AccessToken == "" {
			return "", errors.New("AccessToken not found")
		}

		return adResponse.AccessToken, nil
	}
}

type ADResponse struct {
	AccessToken string `json:"access_token"`
}
