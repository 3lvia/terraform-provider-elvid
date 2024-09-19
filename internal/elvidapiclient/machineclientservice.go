package elvidapiclient

import (
	"encoding/json"
	"fmt"
	"io"
)

func CreateMachineClient(elvidAuthority string, accessTokenAD string, machineClientInput *MachineClientDto) (*MachineClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClientInput)

	response, err := PostRequest(apiUrl, accessTokenAD, machineClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var machineClient MachineClientDto
	err = json.Unmarshal(data, &machineClient)
	if err != nil {
		return nil, err
	}

	return &machineClient, nil
}

func ReadMachineClient(elvidAuthority string, accessTokenAD string, id string) (*MachineClientDto, error) {
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

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var machineClient MachineClientDto
	err = json.Unmarshal(data, &machineClient)
	if err != nil {
		return nil, err
	}

	return &machineClient, nil
}

func UpdateMachineClient(elvidAuthority string, accessTokenAD string, machineClient *MachineClientDto) (*MachineClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClient)

	response, err := PatchRequest(apiUrl, accessTokenAD, machineClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var machineClientResponse MachineClientDto
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
