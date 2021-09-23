package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func CreateUserClient(elvidAuthority string, accessTokenAD string, userClient *UserClient) (*UserClient, error) {
	apiUrl := fmt.Sprintf("%s/api/userclient", elvidAuthority)
	userClientAsJson, _ := json.Marshal(userClient)
	response, err := PostRequest(apiUrl, accessTokenAD, userClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var createdUserClient UserClient
	err = json.Unmarshal(data, &createdUserClient)
	if err != nil {
		return nil, err
	}

	return &createdUserClient, nil
}

func ReadUserClient(elvidAuthority string, accessTokenAD string, id string) (*UserClient, error) {
	apiUrl := fmt.Sprintf("%s/api/userclient/%s", elvidAuthority, id)

	response, err := GetRequest(apiUrl, accessTokenAD)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == 404 {
		return nil, nil
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var userClient UserClient
	err = json.Unmarshal(data, &userClient)
	if err != nil {
		return nil, err
	}

	return &userClient, nil
}

func UpdateUserClient(elvidAuthority string, accessTokenAD string, userClient *UserClient) error {
	apiUrl := fmt.Sprintf("%s/api/userclient", elvidAuthority)

	userClientAsJson, _ := json.Marshal(userClient)
	response, err := PatchRequest(apiUrl, accessTokenAD, userClientAsJson)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}

func DeleteUserClient(elvidAuthority string, accessTokenAD string, id string) error {
	apiUrl := fmt.Sprintf("%s/api/userclient/%s", elvidAuthority, id)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}

type UserClient struct {
	Id                               int      `json:"Id"`
	ClientId                         string   `json:"ClientId"`
	ClientName                       string   `json:"ClientName"`
	Scopes                           []string `json:"Scopes"`
	Domains                          []string `json:"Domains"`
	RedirectUriPaths                 []string `json:"RedirectUriPaths"`
	PostLogoutRedirectUriPaths       []string `json:"PostLogoutRedirectUriPaths"`
	BankIDLoginEnabled               bool     `json:"BankIDLoginEnabled"`
	LocalLoginEnabled                bool     `json:"LocalLoginEnabled"`
	ElviaADLoginEnabled              bool     `json:"ElviaADLoginEnabled"`
	TestUserLoginEnabled             bool     `json:"TestUserLoginEnabled"`
	RequireClientSecret              bool     `json:"RequireClientSecret"`
	AccessTokenLifetime              int      `json:"AccessTokenLifetime"`
	AlwaysIncludeUserClaimsInIdToken bool     `json:"AlwaysIncludeUserClaimsInIdToken"`
	ClientNameLanguageKey            string   `json:"ClientNameLanguageKey"`
	AllowUseOfRefreshTokens          bool     `json:"AllowUseOfRefreshTokens"`
	OneTimeUsageForRefreshTokens     bool     `json:"OneTimeUsageForRefreshTokens"`
	RefreshTokensLifeTime            int      `json:"RefreshTokensLifeTime"`
}
