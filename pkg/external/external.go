package external

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
)

const (
	PostMethod = "POST"
	GetMethod  = "GET"
)

type ExternalService struct {
	CustomHTTPRequest *models.CustomHTTPRequest
}

func (x *ExternalService) PostRequest(url string, request any) (*models.BroadcastResponse, error) {
	resp, err := SendWithCustomRequest(PostMethod, x.CustomHTTPRequest, url, request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &models.BroadcastResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (x *ExternalService) GetRequest(url string) (*models.BroadcastResponse, error) {
	resp, err := SendWithCustomRequest(PostMethod, x.CustomHTTPRequest, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &models.BroadcastResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func SendWithCustomRequest(method string, customRequest *models.CustomHTTPRequest, url string, requestBody any) (*http.Response, error) {
	var err error

	var req *http.Request
	req, err = http.NewRequestWithContext(customRequest.Context, method, url, nil)
	if method == PostMethod && requestBody != nil {
		jsonPayload := &bytes.Buffer{}
		err = json.NewEncoder(jsonPayload).Encode(requestBody)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequestWithContext(customRequest.Context, method, url, jsonPayload)
	}
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	timeout := customRequest.Timeout * time.Second
	client.Timeout = timeout

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
