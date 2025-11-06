package models

type RetryRequest struct {
	RetryAttempt  int
	RetryDuration int
}

type RetryMonitorRequest struct {
	Url          string
	Status       string
	RetryRequest RetryRequest
}
