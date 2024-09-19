package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func CreateOrUpdateApiScope(elvidAuthority string, accessTokenAD string, apiScopeDto *ApiScopeDto) (*ApiScopeDto, error) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, _ := json.Marshal(apiScopeDto)

	response, err := PostRequest(apiUrl, accessTokenAD, apiScopeAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScopeDto
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		return nil, err
	}

	return &apiScope, nil
}

func ReadApiScope(elvidAuthority string, accessTokenAD string, name string, diags diag.Diagnostics) (*ApiScopeDto, diag.Diagnostics) {
	if name == "" {
		diags.AddError("No name provided in ReadApiScope", "")
		return nil, diags
	}

	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, name)

	diags.AddWarning("Calling ApiScope GET in ReadApiScope", "API url = "+apiUrl)

	response, err := GetRequest(apiUrl, accessTokenAD)

	if err != nil {
		diags.AddError("Error from ApiScope GET in ReadApiScope", err.Error())
		return nil, diags
	}

	if response.StatusCode == 404 {
		return nil, diags
	}

	if response.StatusCode != 200 {
		diags.AddError("ApiScope GET returned http error code in ReadApiScope", ElvidErrorResponse(response, apiUrl).Error())
		return nil, diags
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScopeDto
	err = json.Unmarshal(data, &apiScope)

	if err != nil {
		diags.AddError("Could not parse ApiScope GET response as JSON in ReadApiScope", err.Error())
		return nil, diags
	}

	return &apiScope, diags
}

func DeleteApiScope(elvidAuthority string, accessTokenAD string, apiScopeName string) diag.Diagnostics {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, apiScopeName)

	var diags diag.Diagnostics
	diags.AddWarning("Calling ApiScope DELETE in DeleteApiScope", "API url = "+apiUrl)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		diags.AddError("Error from ApiScope DELETE in DeleteApiScope", err.Error())
		return diags
	}

	if response.StatusCode != 200 {
		diags.AddError("ApiScope DELETE returned http error code in DeleteApiScope", ElvidErrorResponse(response, apiUrl).Error())
		return diags
	}

	return diags
}
