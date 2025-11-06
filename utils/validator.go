package utils

import (
	"errors"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
)

func ValidateBroadcastRequest(request *models.BroadcastRequest) error {
	if request == nil {
		return errors.New("request is required")
	}

	if request.Symbol == "" {
		return errors.New("symbol is required")
	}

	if request.Price == 0 {
		return errors.New("price should be greater than 0")
	}

	if request.Timestamp == 0 {
		return errors.New("timestamp should be greater than 0")
	}

	return nil
}
