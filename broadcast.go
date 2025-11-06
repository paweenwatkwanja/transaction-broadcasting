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
	return &BroadcastService{}
}

func (b *BroadcastService) BroadcastTransaction(url string, request models.BroadcastRequest) (string, error) {
	txHash := ""

	err := utils.ValidateBroadcastRequest(request)
	if err != nil {
		return txHash, err
	}

	response, err := b.externalService.PostRequest(url, request)
	if err != nil {
		return txHash, err
	}
	txHash = response.TxHash

	return txHash, nil
}

func (b *BroadcastService) MonitorTransaction(url string) (string, error) {
	txStatus := ""

	response, err := b.externalService.GetRequest(url)
	if err != nil {
		return txStatus, err
	}
	txStatus = response.TxStatus

	return txStatus, nil
}

func (b *BroadcastService) HandleStatus(url string, txStatus string) error {
	switch txStatus {
	case confirmedStatus:
		return nil
	case failedStatus:
		return errors.New("broadcast failed")
	case pendingStatus:
		return b.retryMonitorTransaction(url, confirmedStatus, b.retryRequest)
	case dneStatus:
		return errors.New("item not exist")
	default:
		return errors.New("status not exist")
	}
}

func (b *BroadcastService) retryMonitorTransaction(url string, status string, retryRequest *models.RetryRequest) error {
	var retryAttempt uint = 3
	if retryRequest.RetryAttempt != 0 {
		retryAttempt = retryRequest.RetryAttempt
	}
	var retryDuration time.Duration = 5
	if retryRequest.RetryAttempt != 0 {
		retryAttempt = retryRequest.RetryAttempt
	}

	for i := range retryAttempt {
		response, err := b.externalService.GetRequest(url)
		if err != nil {
			return err
		}
		txStatus := response.TxStatus

		if txStatus == status {
			break
		} else if txStatus == failedStatus {
			return errors.New("broadcast failed")
		} else if txStatus == dneStatus {
			return errors.New("item not found")
		}

		if i < retryAttempt {
			fmt.Printf("Attempt %v failed. Retrying in %v seconds\n", i+1, retryDuration)
			time.Sleep(retryDuration * time.Second)
		}
	}

	return errors.New("could not get confirmed status")
}

func (b *BroadcastService) WithRetryAttempt(retryAttempt uint) {
	b.retryRequest.RetryAttempt = retryAttempt
}

func (b *BroadcastService) WithRetryDuration(retryDuration time.Duration) {
	b.retryRequest.RetryDuration = retryDuration
}

func (b *BroadcastService) WithCustomizedHTTPRequest(customHTTPRequest *models.CustomHTTPRequest) {
	b.externalService.CustomHTTPRequest = customHTTPRequest
}
