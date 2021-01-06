package elvidapiclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 15 * time.Second}

func PostRequest(url string, accessToken string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	return myClient.Do(req)
}

func PatchRequest(url string, accessToken string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	return myClient.Do(req)
}

func GetRequest(url string, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return myClient.Do(req)
}

func DeleteRequest(url string, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return myClient.Do(req)
}

func ElvidErrorResponse(response *http.Response, url string) error {
	data, _ := ioutil.ReadAll(response.Body)
	return fmt.Errorf("ElvID retured StatusCode %v for (%s), message: %s", response.StatusCode, url, data)
}
