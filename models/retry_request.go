package models

import "time"

type RetryRequest struct {
	RetryAttempt  int
	RetryDuration time.Duration
}

type RetryMonitorRequest struct {
	Url          string
	Status       string
	RetryRequest RetryRequest
}
