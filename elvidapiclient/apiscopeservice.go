package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func CreateOrUpdateApiScope(elvidAuthority string, accessTokenAD string, apiScopeInput *ApiScope) (*ApiScope, error) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, _ := json.Marshal(apiScopeInput)

	response, err := PostRequest(apiUrl, accessTokenAD, apiScopeAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		return nil, err
	}

	return &apiScope, nil
}

func ReadApiScope(elvidAuthority string, accessTokenAD string, id string) (*ApiScope, error) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, id)

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

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		return nil, err
	}

	return &apiScope, nil
}

func DeleteApiScope(elvidAuthority string, accessTokenAD string, apiScopeName string) error {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, apiScopeName)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}

type ApiScope struct {
	Name                string   `json:"Name"`
	Description         string   `json:"Description"`
	UserClaims          []string `json:"UserClaims"`
	AllowMachineClients bool     `json:"AllowMachineClients"`
	AllowUserClients    bool     `json:"AllowUserClients"`
}
