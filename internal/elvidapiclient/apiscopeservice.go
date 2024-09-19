package elvidapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func CreateOrUpdateApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, apiScopeInput *ApiScope) (*ApiScope, diag.Diagnostics) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, _ := json.Marshal(apiScopeInput)

	var diags diag.Diagnostics
	diags.AddWarning("Calling ApiScope POST in CreateOrUpdateApiScope", "API url = "+apiUrl+", Api scope JSON = "+string(apiScopeAsJson))

	response, err := PostRequest(apiUrl, accessTokenAD, apiScopeAsJson)

	if err != nil {
		diags.AddError("ApiScope POST error in CreateOrUpdateApiScope", err.Error())
		return nil, diags
	}

	if response.StatusCode != 200 {
		diags.AddError("ApiScope POST returned http error code in CreateOrUpdateApiScope", ElvidErrorResponse(response, apiUrl).Error())
		return nil, diags
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		diags.AddError("Could not parse ApiScope POST response as JSON in CreateOrUpdateApiScope", err.Error())
		return nil, diags
	}

	return &apiScope, diags
}

func ReadApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, name string) (*ApiScope, diag.Diagnostics) {
	var diags diag.Diagnostics

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

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)

	if err != nil {
		diags.AddError("Could not parse ApiScope GET response as JSON in ReadApiScope", err.Error())
		return nil, diags
	}

	return &apiScope, diags
}

func DeleteApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, apiScopeName string) diag.Diagnostics {
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

// TODO: Don't use terraform type in elvidapiclient
type ApiScope struct {
	Name                types.String `json:"Name"`
	Description         types.String `json:"Description"`
	UserClaims          types.List   `json:"UserClaims"`
	AllowMachineClients types.Bool   `json:"AllowMachineClients"`
	AllowUserClients    types.Bool   `json:"AllowUserClients"`
}
