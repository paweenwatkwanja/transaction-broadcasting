package broadcast

import (
	"errors"
	"fmt"
	"time"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
	"github.com/paweenwatkwanja/transaction-broadcasting/pkg/external"
	"github.com/paweenwatkwanja/transaction-broadcasting/utils"
)

type BroadcastService struct{}

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

	responseBody, err := external.PostRequest(url, request)
	if err != nil {
		return txHash, err
	}
	txHash = string(responseBody)

	return txHash, nil
}

func (b *BroadcastService) MonitorTransaction(url string) (string, error) {
	txStatus := ""

	responseBody, err := external.GetRequest(url)
	if err != nil {
		return txStatus, err
	}
	txStatus = string(responseBody)

	return txStatus, nil
}

func (b *BroadcastService) HandleStatus(url string, txStatus string) error {
	switch txStatus {
	case confirmedStatus:
		return nil
	case failedStatus:
		return errors.New("broadcast failed")
	case pendingStatus:
		return retryMonitorTransaction(url, confirmedStatus)
	case dneStatus:
		return errors.New("item not found")
	default:
		return errors.New("status not exist")
	}
}

func retryMonitorTransaction(url string, status string) error {
	retryAttempt := 3
	duration := 5

	for i := range retryAttempt {
		responseBody, err := external.GetRequest(url)
		if err != nil {
			return err
		}
		txStatus := string(responseBody)

		if txStatus == status {
			break
		} else if txStatus == failedStatus {
			return errors.New("broadcast failed")
		} else if txStatus == dneStatus {
			return errors.New("status not exist")
		}

		if i < retryAttempt {
			fmt.Printf("attempt %v failed. Retrying in %v seconds", i, duration)
			time.Sleep(time.Duration(duration) * time.Second)
		}
	}

	return errors.New("could not get confirmed status")
}
