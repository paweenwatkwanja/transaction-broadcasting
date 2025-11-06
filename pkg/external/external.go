package external

import (
	"time"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
	"resty.dev/v3"
)

type ExternalService struct {
	CustomHTTPRequest *models.CustomHTTPRequest
}

func (x *ExternalService) Post(url string, request *models.BroadcastRequest) (*models.BroadcastResponse, error) {
	client := initClient(x.CustomHTTPRequest)
	defer client.Close()

	resp, err := client.R().
		SetBody(request).
		SetResult(&models.BroadcastResponse{}).
		Post(url)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*models.BroadcastResponse), nil
}

func (x *ExternalService) Get(url string) (*models.BroadcastResponse, error) {
	client := initClient(x.CustomHTTPRequest)
	defer client.Close()

	resp, err := client.R().
		SetResult(&models.BroadcastResponse{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*models.BroadcastResponse), nil
}

func initClient(request *models.CustomHTTPRequest) *resty.Client {
	client := resty.New()
	if request == nil {
		return client
	}

	if request.RetryAttempt != 0 {
		client.SetRetryCount(request.RetryAttempt)
	}

	if request.RetryDuration != 0 {
		client.SetRetryWaitTime(time.Duration(request.RetryDuration) * time.Second).
			SetRetryMaxWaitTime(time.Duration(request.RetryDuration) * time.Second)
	}

	if request.Timeout != 0 {
		client.SetTimeout(time.Duration(request.RetryDuration) * time.Second)
	}

	return client
}
