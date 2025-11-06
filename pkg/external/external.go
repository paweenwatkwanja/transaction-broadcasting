package external

import (
	"fmt"
	"time"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
	"resty.dev/v3"
)

type ExternalService struct {
	CustomHTTPRequest *models.CustomHTTPRequest
}

func (x *ExternalService) Post(url string, request *models.BroadcastRequest) (*models.BroadcastResponse, error) {
	fmt.Println("Post method is called")
	client := initClient(x.CustomHTTPRequest)
	defer client.Close()

	resp, err := client.R().
		SetBody(request).
		SetResult(&models.BroadcastResponse{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	return resp.Result().(*models.BroadcastResponse), nil
}

func (x *ExternalService) Get(url string) (*models.BroadcastResponse, error) {
	fmt.Println("Get method is called")
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
	fmt.Println("Init Client")
	fmt.Println(request)
	client := resty.New()
	if request == nil {
		return client
	}

	if request.RetryAttempt != 0 {
		fmt.Printf("RetryAttempt is %v", request.RetryAttempt)
		client.SetRetryCount(request.RetryAttempt)
	}

	if request.RetryDuration != 0 {
		fmt.Printf("RetryDuration is %v", request.RetryDuration)
		client.SetRetryWaitTime(request.RetryDuration * time.Second).
			SetRetryMaxWaitTime(request.RetryDuration * time.Second)
	}

	if request.Timeout != 0 {
		fmt.Printf("Timeout is %v", request.Timeout)
		client.SetTimeout(request.Timeout)
	}

	return client
}
