package external

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
)

func PostRequest(url string, request any) (*models.BroadcastResponse, error) {
	contentType := "application/json"
	requestBody := &bytes.Buffer{}
	err := json.NewEncoder(requestBody).Encode(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, contentType, requestBody)
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

func GetRequest(url string) (*models.BroadcastResponse, error) {
	resp, err := http.Get(url)
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
