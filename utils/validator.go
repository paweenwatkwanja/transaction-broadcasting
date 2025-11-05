package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/paweenwatkwanja/transaction-broadcasting/models"
)

func ValidateBroadcastRequest(request models.BroadcastRequest) error {
	validator := validator.New()
	err := validator.Struct(request)
	if err != nil {
		return errors.New("request is required")
	}

	return nil
}
