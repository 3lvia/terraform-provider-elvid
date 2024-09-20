package elvidapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func CreateUserClient(ctx context.Context, elvidAuthority string, accessTokenAD string, userClient *UserClientDto) (*UserClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/userclient", elvidAuthority)
	userClientAsJson, _ := json.Marshal(userClient)
	response, err := PostRequest(ctx, apiUrl, accessTokenAD, userClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var createdUserClient UserClientDto
	err = json.Unmarshal(data, &createdUserClient)
	if err != nil {
		return nil, err
	}

	return &createdUserClient, nil
}

func ReadUserClient(ctx context.Context, elvidAuthority string, accessTokenAD string, id string) (*UserClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/userclient/%s", elvidAuthority, id)

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

	var userClient UserClientDto
	err = json.Unmarshal(data, &userClient)
	if err != nil {
		return nil, err
	}

	return &userClient, nil
}

func UpdateUserClient(ctx context.Context, elvidAuthority string, accessTokenAD string, userClient *UserClientDto) error {
	apiUrl := fmt.Sprintf("%s/api/userclient", elvidAuthority)

	userClientAsJson, _ := json.Marshal(userClient)
	response, err := PatchRequest(ctx, apiUrl, accessTokenAD, userClientAsJson)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}

func DeleteUserClient(ctx context.Context, elvidAuthority string, accessTokenAD string, id string) error {
	apiUrl := fmt.Sprintf("%s/api/userclient/%s", elvidAuthority, id)

	response, err := DeleteRequest(ctx, apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
