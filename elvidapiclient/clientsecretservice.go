package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

func CreateClientSecret(elvidAuthority string, accessTokenAD string, clientId string) (*ClientSecret, error) {
	apiUrl := fmt.Sprintf("%s/api/clientsecret", elvidAuthority)

	clientIdAsInt, _ := strconv.Atoi(clientId)
	values := map[string]interface{}{
		"ClientId": clientIdAsInt,
	}

	jsonValue, err := json.Marshal(values)
	response, err := PostRequest(apiUrl, accessTokenAD, jsonValue)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var clientSecret ClientSecret
	err = json.Unmarshal(data, &clientSecret)
	if err != nil {
		return nil, err
	}

	return &clientSecret, nil
}

func ReadClientSecret(elvidAuthority string, accessTokenAD string, clientId string, clientSecretId string) (*ClientSecret, error) {
	apiUrl := fmt.Sprintf("%s/api/clientsecret/%s/%s", elvidAuthority, clientId, clientSecretId)

	response, err := GetRequest(apiUrl, accessTokenAD)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
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

	var clientSecret ClientSecret
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

type ClientSecret struct {
	Id                    int    `json:"Id"`
	Value                 string `json:"SecretValue"`
	HashedValueStartsWith string `json:"HashedValueStartsWith"`
}
