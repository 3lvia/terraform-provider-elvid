package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func CreateMachineClient(elvidAuthority string, accessTokenAD string, machineClientInput *MachineClient) (*MachineClient, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClientInput)

	response, err := PostRequest(apiUrl, accessTokenAD, machineClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var machineClient MachineClient
	err = json.Unmarshal(data, &machineClient)
	if err != nil {
		return nil, err
	}

	return &machineClient, nil
}

func ReadMachineClient(elvidAuthority string, accessTokenAD string, id string) (*MachineClient, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient/%s", elvidAuthority, id)

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

	var machineClient MachineClient
	err = json.Unmarshal(data, &machineClient)
	if err != nil {
		return nil, err
	}

	return &machineClient, nil
}

func UpdateMachineClient(elvidAuthority string, accessTokenAD string, machineClient *MachineClient) (*MachineClient, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClient)

	response, err := PatchRequest(apiUrl, accessTokenAD, machineClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var machineClientResponse MachineClient
	err = json.Unmarshal(data, &machineClient)
	if err != nil {
		return nil, err
	}

	return &machineClientResponse, nil
}

func DeleteMachineClient(elvidAuthority string, accessTokenAD string, id string) error {
	apiUrl := fmt.Sprintf("%s/api/machineclient/%s", elvidAuthority, id)

	response, err := DeleteRequest(apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}

type MachineClient struct {
	Id                   int           `json:"Id"`
	ClientId             string        `json:"ClientId"`
	ClientName           string        `json:"ClientName"`
	TestUserLoginEnabled bool          `json:"TestUserLoginEnabled"`
	IsDelegationClient   bool          `json:"IsDelegationClient"`
	AccessTokenLifeTime  int           `json:"AccessTokenLifeTime"`
	Scopes               []string      `json:"Scopes"`
	ClientClaims         []ClientClaim `json:"ClientClaims"`
}

type ClientClaim struct {
	Type   string   `json:"Type"`
	Values []string `json:"Values"`
}
