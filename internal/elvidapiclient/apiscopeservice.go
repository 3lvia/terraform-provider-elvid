package elvidapiclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

func CreateOrUpdateApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, apiScopeDto *ApiScopeDto) (*ApiScopeDto, error) {
	apiUrl := fmt.Sprintf("%s/api/ApiScope", elvidAuthority)

	apiScopeAsJson, err := json.Marshal(apiScopeDto)
	if err != nil {
		return nil, err
	}

	response, err := PostRequest(ctx, apiUrl, accessTokenAD, apiScopeAsJson)
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

// ValidateNewApiScope dry-runs the create-time validation server-side without
// persisting anything. Used by ApiScopeResource.ModifyPlan so naming-convention
// and AD-group-existence errors surface in the speculative plan on a PR instead
// of mid-apply on merge. See ADR 2026-05-CORE-2650 (H2) in the elvid repo.
func ValidateNewApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, apiScopeDto *ApiScopeDto) error {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/validate", elvidAuthority)

	apiScopeAsJson, err := json.Marshal(apiScopeDto)
	if err != nil {
		return err
	}

	response, err := PostRequest(ctx, apiUrl, accessTokenAD, apiScopeAsJson)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return nil
	}
	return ElvidErrorResponse(response, apiUrl)
}

func ReadApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, name string) (*ApiScopeDto, error) {
	if name == "" {
		return nil, errors.New("no name provided in ReadApiScope")
	}

	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, name)

	response, err := GetRequest(ctx, apiUrl, accessTokenAD)

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

func DeleteApiScope(ctx context.Context, elvidAuthority string, accessTokenAD string, apiScopeName string) error {
	apiUrl := fmt.Sprintf("%s/api/ApiScope/%s", elvidAuthority, apiScopeName)

	response, err := DeleteRequest(ctx, apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
