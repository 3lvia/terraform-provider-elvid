package elvidapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func CreateClientSecret(ctx context.Context, elvidAuthority string, accessTokenAD string, clientId string) (*ClientSecretDto, error) {
	apiUrl := fmt.Sprintf("%s/api/clientsecret", elvidAuthority)

	clientIdAsInt, _ := strconv.Atoi(clientId)
	values := map[string]interface{}{
		"ClientId": clientIdAsInt,
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	response, err := PostRequest(ctx, apiUrl, accessTokenAD, jsonValue)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var clientSecret ClientSecretDto
	err = json.Unmarshal(data, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &clientSecret, nil
}

func ReadClientSecret(ctx context.Context, elvidAuthority string, accessTokenAD string, clientId string, clientSecretId string) (*ClientSecretDto, error) {
	apiUrl := fmt.Sprintf("%s/api/clientsecret/%s/%s", elvidAuthority, clientId, clientSecretId)

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

	var clientSecret ClientSecretDto
	err = json.Unmarshal(data, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &clientSecret, nil
}

func DeleteClientSecret(ctx context.Context, elvidAuthority string, accessTokenAD string, clientId string, clientSecretId string) error {
	apiUrl := fmt.Sprintf("%s/api/clientsecret/%s/%s", elvidAuthority, clientId, clientSecretId)

	response, err := DeleteRequest(ctx, apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
