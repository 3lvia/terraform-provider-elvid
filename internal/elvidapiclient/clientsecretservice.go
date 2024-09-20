package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

func CreateClientSecret(elvidAuthority string, accessTokenAD string, clientId string) (*ClientSecretDto, error) {
	ioutil.WriteFile("logs/CreateClientSecret.text", []byte("someString"), 0644)
	apiUrl := fmt.Sprintf("%s/api/clientsecret", elvidAuthority)

	clientIdAsInt, _ := strconv.Atoi(clientId)
	values := map[string]interface{}{
		"ClientId": clientIdAsInt,
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	response, err := PostRequest(apiUrl, accessTokenAD, jsonValue)

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

func ReadClientSecret(elvidAuthority string, accessTokenAD string, clientId string, clientSecretId string) (*ClientSecretDto, error) {
	apiUrl := fmt.Sprintf("%s/api/clientsecret/%s/%s", elvidAuthority, clientId, clientSecretId)

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

	var clientSecret ClientSecretDto
	err = json.Unmarshal(data, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &clientSecret, nil
}

func DeleteClientSecret(elvidAuthority string, accessTokenAD string, clientId string, clientSecretId string) error {
	apiUrl := fmt.Sprintf("%s/api/clientsecret/%s/%s", elvidAuthority, clientId, clientSecretId)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
