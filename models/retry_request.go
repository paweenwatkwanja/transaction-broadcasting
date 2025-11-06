package models

import "time"

type RetryRequest struct {
	RetryAttempt  uint
	RetryDuration time.Duration
}
