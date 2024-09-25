package elvidapiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func CreateMachineClient(ctx context.Context, elvidAuthority string, accessTokenAD string, machineClientInput *MachineClientDto) (*MachineClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClientInput)

	response, err := PostRequest(ctx, apiUrl, accessTokenAD, machineClientAsJson)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, readErr := io.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, readErr
	}

	var machineClientResponseDto MachineClientDto
	err = json.Unmarshal(data, &machineClientResponseDto)
	if err != nil {
		return nil, err
	}

	return &machineClientResponseDto, nil
}

func ReadMachineClient(ctx context.Context, elvidAuthority string, accessTokenAD string, id string) (*MachineClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient/%s", elvidAuthority, id)

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

	data, readErr := io.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		return nil, readErr
	}

	var machineClientResponseDto MachineClientDto
	err = json.Unmarshal(data, &machineClientResponseDto)
	if err != nil {
		return nil, err
	}

	return &machineClientResponseDto, nil
}

func UpdateMachineClient(ctx context.Context, elvidAuthority string, accessTokenAD string, machineClient *MachineClientDto) (*MachineClientDto, error) {
	apiUrl := fmt.Sprintf("%s/api/machineclient", elvidAuthority)

	machineClientAsJson, _ := json.Marshal(machineClient)

	response, err := PatchRequest(ctx, apiUrl, accessTokenAD, machineClientAsJson)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, ElvidErrorResponse(response, apiUrl)
	}

	data, readErr := io.ReadAll(response.Body)
	defer response.Body.Close()
	if readErr != nil {
		return nil, readErr
	}

	var machineClientResponseDto MachineClientDto
	err = json.Unmarshal(data, &machineClientResponseDto)
	if err != nil {
		return nil, err
	}

	return &machineClientResponseDto, nil
}

func DeleteMachineClient(ctx context.Context, elvidAuthority string, accessTokenAD string, id string) error {
	apiUrl := fmt.Sprintf("%s/api/machineclient/%s", elvidAuthority, id)

	response, err := DeleteRequest(ctx, apiUrl, accessTokenAD)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return ElvidErrorResponse(response, apiUrl)
	}

	return nil
}
