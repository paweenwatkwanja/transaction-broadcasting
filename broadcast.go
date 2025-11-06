package broadcast

import (
	"errors"
	"fmt"
	"time"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
	"github.com/paweenwatkwanja/transaction-broadcasting/pkg/external"
	"github.com/paweenwatkwanja/transaction-broadcasting/utils"
)

type BroadcastService struct {
	retryRequest    *models.RetryRequest
	externalService external.ExternalService
}

const (
	confirmedStatus = "COMFIRMED"
	failedStatus    = "FAILED"
	pendingStatus   = "PENDING"
	dneStatus       = "DNE"
)

func NewBroadcastService() *BroadcastService {
	return &BroadcastService{
		retryRequest: &models.RetryRequest{},
	}
}

func (b *BroadcastService) BroadcastTransaction(url string, request *models.BroadcastRequest) (string, error) {
	txHash := ""
	fmt.Println("broadcasting")
	err := utils.ValidateBroadcastRequest(request)
	if err != nil {
		return txHash, err
	}
	fmt.Println("request validated")

	response, err := b.externalService.Post(url, request)
	if err != nil {
		return txHash, err
	}
	fmt.Println("got response from Post")
	txHash = response.TxHash
	fmt.Printf("txHash is %v\n", response.TxHash)
	return txHash, nil
}

func (b *BroadcastService) MonitorTransaction(url string) (string, error) {
	txStatus := ""
	fmt.Println("monitoring")
	response, err := b.externalService.Get(url)
	if err != nil {
		return txStatus, err
	}
	fmt.Println("got response from Get")
	txStatus = response.TxStatus
	fmt.Printf("txStatus is %v\n", txStatus)
	return txStatus, nil
}

func (b *BroadcastService) HandleStatus(url string, txStatus string) error {
	fmt.Println("handling")
	switch txStatus {
	case confirmedStatus:
		fmt.Println("case confirmed")
		return nil
	case failedStatus:
		fmt.Println("case failed")
		return errors.New("broadcast failed")
	case pendingStatus:
		fmt.Println("case pedning")
		retryMonitorRequest := &models.RetryMonitorRequest{
			Url:          url,
			Status:       txStatus,
			RetryRequest: *b.retryRequest,
		}
		return b.retryMonitorTransaction(retryMonitorRequest)
	case dneStatus:
		fmt.Println("case dns")
		return errors.New("item not exist")
	default:
		fmt.Println("status not exist")
		return errors.New("status not exist")
	}
}

func (b *BroadcastService) retryMonitorTransaction(retryMonitorRequest *models.RetryMonitorRequest) error {
	fmt.Println("retrying")
	retryRequest := &retryMonitorRequest.RetryRequest

	var retryAttempt int = 3
	if retryRequest.RetryAttempt != 0 {
		retryAttempt = retryRequest.RetryAttempt
	}
	var retryDuration time.Duration = 5
	if retryRequest.RetryDuration != 0 {
		retryDuration = retryRequest.RetryDuration
	}

	for i := range retryAttempt {
		fmt.Println(i)
		fmt.Println(retryAttempt)
		fmt.Println(retryDuration)
		response, err := b.externalService.Get(retryMonitorRequest.Url)
		fmt.Println(response)
		if err != nil {
			return err
		}
		txStatus := response.TxStatus

		if txStatus == confirmedStatus {
			break
		} else if txStatus == failedStatus {
			return errors.New("broadcast failed")
		} else if txStatus == dneStatus {
			return errors.New("item not found")
		}

		if i < retryAttempt {
			fmt.Printf("Attempt %v failed. Retrying in %v seconds\n", i+1, retryDuration.Seconds())
			time.Sleep(retryDuration * time.Second)
		}
	}

	return errors.New("could not get confirmed status")
}

func (b *BroadcastService) WithRetryRequest(retryRequest *models.RetryRequest) *BroadcastService {
	b.retryRequest = retryRequest
	fmt.Printf("retryRequest : %v", b.retryRequest)
	return b
}

func (b *BroadcastService) WithCustomHTTPRequest(customHTTPRequest *models.CustomHTTPRequest) *BroadcastService {
	b.externalService.CustomHTTPRequest = &models.CustomHTTPRequest{}
	b.externalService.CustomHTTPRequest = customHTTPRequest
	fmt.Printf("CustomHTTPRequest : %v", b.externalService.CustomHTTPRequest)
	return b
}
