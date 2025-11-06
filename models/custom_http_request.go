package models

import (
	"time"
)

type CustomHTTPRequest struct {
	RetryAttempt  int
	RetryDuration time.Duration
	Timeout       time.Duration
}
