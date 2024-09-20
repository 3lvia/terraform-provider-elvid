package elvidapiclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func CreateOrUpdateApiScope(elvidAuthority string, accessTokenAD string, apiScopeDto *ApiScopeDto) (*ApiScopeDto, error) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, err := json.Marshal(apiScopeDto)
	if err != nil {
		return nil, err
	}

	response, err := PostRequest(apiUrl, accessTokenAD, apiScopeAsJson)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScopeDto
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		return nil, err
	}

	return &apiScope, nil
}

func ReadApiScope(elvidAuthority string, accessTokenAD string, name string) (*ApiScopeDto, error) {
	if name == "" {
		return nil, errors.New("no name provided in ReadApiScope")
	}

	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, name)

	// diags.AddWarning("Calling ApiScope GET in ReadApiScope", "API url = "+apiUrl)

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

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScopeDto
	err = json.Unmarshal(data, &apiScope)

	if err != nil {
		return nil, err
	}

	return &apiScope, nil
}

func DeleteApiScope(elvidAuthority string, accessTokenAD string, apiScopeName string) error {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, apiScopeName)

	// diags.AddWarning("Calling ApiScope DELETE in DeleteApiScope", "API url = "+apiUrl)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
