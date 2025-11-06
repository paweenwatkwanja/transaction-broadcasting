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
	confirmedStatus = "CONFIRMED"
	failedStatus    = "FAILED"
	pendingStatus   = "PENDING"
	dneStatus       = "DNE"
)

func NewBroadcastService() *BroadcastService {
	return &BroadcastService{}
}

func (b *BroadcastService) BroadcastTransaction(url string, request *models.BroadcastRequest) (string, error) {
	txHash := ""

	err := utils.ValidateBroadcastRequest(request)
	if err != nil {
		return txHash, err
	}

	response, err := b.externalService.Post(url, request)
	if err != nil {
		return txHash, err
	}
	txHash = response.TxHash

	return txHash, nil
}

func (b *BroadcastService) MonitorTransaction(url string) (string, error) {
	txStatus := ""

	response, err := b.externalService.Get(url)
	if err != nil {
		return txStatus, err
	}
	txStatus = response.TxStatus

	return txStatus, nil
}

func (b *BroadcastService) HandleStatus(url string, txStatus string) (string, error) {
	switch txStatus {
	case confirmedStatus:
		return "", nil
	case failedStatus:
		return txStatus, errors.New("broadcast failed")
	case pendingStatus:
		retryMonitorRequest := newRetryMonitorRequest(url, txStatus, b.retryRequest)
		return b.retryMonitorTransaction(retryMonitorRequest)
	case dneStatus:
		return txStatus, errors.New("item not exist")
	default:
		return txStatus, errors.New("status not exist")
	}
}

func (b *BroadcastService) retryMonitorTransaction(retryMonitorRequest *models.RetryMonitorRequest) (string, error) {
	retryRequest := &retryMonitorRequest.RetryRequest
	retryAttempt := getRetryAttemptForBroadcast(retryRequest.RetryAttempt)
	retryDuration := getRetryDurationForBroadcast(retryRequest.RetryDuration)

	var txStatus string
	for i := 0; i < retryAttempt; i++ {
		response, err := b.externalService.Get(retryMonitorRequest.Url)
		if err != nil {
			return "", err
		}
		txStatus = response.TxStatus

		switch txStatus {
		case confirmedStatus:
			return txStatus, nil
		case failedStatus:
			return txStatus, errors.New("broadcast failed")
		case dneStatus:
			return txStatus, errors.New("item not found")
		}

		if i < retryAttempt-1 {
			fmt.Printf("Attempt %v. Status is still PENDING. Retrying in %v seconds\n", i+1, retryDuration)
			time.Sleep(time.Duration(retryDuration) * time.Second)
		} else {
			fmt.Printf("Attempt %v. Status is still PENDING. No more retries left.\n", i+1)
		}
	}

	return txStatus, errors.New("status is still pending")
}

func (b *BroadcastService) WithRetryRequest(retryRequest *models.RetryRequest) *BroadcastService {
	b.retryRequest = retryRequest
	return b
}

func (b *BroadcastService) WithCustomHTTPRequest(customHTTPRequest *models.CustomHTTPRequest) *BroadcastService {
	b.externalService.CustomHTTPRequest = customHTTPRequest
	return b
}

func newRetryMonitorRequest(url string, txStatus string, retryRequest *models.RetryRequest) *models.RetryMonitorRequest {
	retryMonitorRequest := &models.RetryMonitorRequest{
		Url:    url,
		Status: txStatus,
	}
	if retryRequest != nil {
		retryMonitorRequest.RetryRequest = *retryRequest
	}

	return retryMonitorRequest
}

func getRetryAttemptForBroadcast(newRetryAttempt int) int {
	var retryAttempt int = 3
	if newRetryAttempt != 0 {
		retryAttempt = newRetryAttempt
	}
	return retryAttempt
}

func getRetryDurationForBroadcast(newRetryDuration int) int {
	var retryDuration int = 5
	if newRetryDuration != 0 {
		retryDuration = newRetryDuration
	}
	return retryDuration
}
