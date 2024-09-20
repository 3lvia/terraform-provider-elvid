package elvidapiclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var myClient = &http.Client{Timeout: 15 * time.Second}

func PostRequest(ctx context.Context, url string, accessToken string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	return DoRequestAndLog(ctx, req)
}

func PatchRequest(ctx context.Context, url string, accessToken string, jsonValue []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	return DoRequestAndLog(ctx, req)
}

func GetRequest(ctx context.Context, url string, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return DoRequestAndLog(ctx, req)
}

func DeleteRequest(ctx context.Context, url string, accessToken string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	return DoRequestAndLog(ctx, req)
}

func DoRequestAndLog(ctx context.Context, req *http.Request) (*http.Response, error) {
	response, err := myClient.Do(req)
	if err == nil {
		if response.StatusCode < 300 {
			tflog.Info(ctx, "HTTP request: "+req.Method+" "+req.URL.String()+" returned "+response.Status)
		} else {
			tflog.Warn(ctx, "HTTP request: "+req.Method+" "+req.URL.String()+" returned "+response.Status)
		}

	} else {
		tflog.Error(ctx, "HTTP request failed: "+req.Method+" "+req.URL.String())
	}
	return response, err
}

func ElvidErrorResponse(response *http.Response, url string) error {
	data, _ := io.ReadAll(response.Body)
	return fmt.Errorf("ElvID returned StatusCode %v for (%s), message: %s", response.StatusCode, url, data)
}
