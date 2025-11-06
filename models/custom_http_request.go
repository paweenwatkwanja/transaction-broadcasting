package models

import (
	"time"
)

type CustomHTTPRequest struct {
	RetryAttempt  int
	RetryDuration time.Duration
	RetryRequest  RetryRequest
	Timeout       time.Duration
}
