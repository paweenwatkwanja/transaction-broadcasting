package external

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PostRequest(url string, request any) ([]byte, error) {
	contentType := "application/json"
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, contentType, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func GetRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
