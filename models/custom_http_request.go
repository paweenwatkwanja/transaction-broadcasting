package models

import (
	"context"
	"time"
)

type CustomHTTPRequest struct {
	RetryRequest RetryRequest
	Timeout      time.Duration
	Context      context.Context
}
