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
	fmt.Printf("result from Post is %v\n", resp.Result())
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
	fmt.Printf("result from Get is %v\n", resp.Result())
	return resp.Result().(*models.BroadcastResponse), nil
}

func initClient(request *models.CustomHTTPRequest) *resty.Client {
	fmt.Println("Init Client")
	client := resty.New()
	if request == nil {
		return client
	}
	fmt.Println(*request)

	if request.RetryAttempt != 0 {
		test := &request.RetryAttempt
		fmt.Printf("RetryAttempt is %v\n", *test)
		client.SetRetryCount(request.RetryAttempt)
	}

	if request.RetryDuration != 0 {
		test := &request.RetryDuration
		fmt.Printf("RetryDuration is %v\n", *test)
		client.SetRetryWaitTime(request.RetryDuration * time.Second).
			SetRetryMaxWaitTime(request.RetryDuration * time.Second)
	}

	if request.Timeout != 0 {
		test := &request.Timeout
		fmt.Printf("Timeout is %v\n", *test)
		client.SetTimeout(request.Timeout)
	}

	return client
}
