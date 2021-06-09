package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func CreateOrUpdateApiScope(elvidAuthority string, accessTokenAD string, apiScopeInput *ApiScope) (*ApiScope, diag.Diagnostics) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, _ := json.Marshal(apiScopeInput)

	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning, // Add a warning (debug info) that will only be shown if something errors
		Summary:  "Calling ApiScope POST in CreateOrUpdateApiScope",
		Detail:   "API url = " + apiUrl + ", Api scope JSON = " + string(apiScopeAsJson),
	})

	response, err := PostRequest(apiUrl, accessTokenAD, apiScopeAsJson)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ApiScope POST error in CreateOrUpdateApiScope",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	if response.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ApiScope POST returned http error code in CreateOrUpdateApiScope",
			Detail:   ElvidErrorResponse(response, apiUrl).Error(),
		})
		return nil, diags
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse ApiScope POST response as JSON in CreateOrUpdateApiScope",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return &apiScope, nil
}

func ReadApiScope(elvidAuthority string, accessTokenAD string, name string) (*ApiScope, diag.Diagnostics) {
	var diags diag.Diagnostics

	if name == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No name provided in ReadApiScope",
			Detail:   "",
		})
		return nil, diags
	}

	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, name)

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning, // Add a warning (debug info) that will only be shown if something errors
		Summary:  "Calling ApiScope GET in ReadApiScope",
		Detail:   "API url = " + apiUrl,
	})

	response, err := GetRequest(apiUrl, accessTokenAD)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error from ApiScope GET in ReadApiScope",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	if response.StatusCode == 404 {
		return nil, nil
	}

	if response.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ApiScope GET returned http error code in ReadApiScope",
			Detail:   ElvidErrorResponse(response, apiUrl).Error(),
		})
		return nil, diags
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var apiScope ApiScope
	err = json.Unmarshal(data, &apiScope)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Could not parse ApiScope GET response as JSON in ReadApiScope",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return &apiScope, nil
}

func DeleteApiScope(elvidAuthority string, accessTokenAD string, apiScopeName string) diag.Diagnostics {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, apiScopeName)

	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning, // Add a warning (debug info) that will only be shown if something errors
		Summary:  "Calling ApiScope DELETE in DeleteApiScope",
		Detail:   "API url = " + apiUrl,
	})

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error from ApiScope DELETE in DeleteApiScope",
			Detail:   err.Error(),
		})
		return diags
	}

	if response.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ApiScope DELETE returned http error code in DeleteApiScope",
			Detail:   ElvidErrorResponse(response, apiUrl).Error(),
		})
		return diags
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
