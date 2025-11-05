package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/paweenwatkwanja/transaction-broadcasting/models"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	testCaseName string
	request      models.BroadcastRequest
	err          error
}

func TestBroadcastRequest(t *testing.T) {
	testCases := []testCase{
		{
			testCaseName: "case request not nil",
			request: models.BroadcastRequest{
				Symbol:    "BTC",
				Price:     100000,
				Timestamp: 1678912345,
			},
			err: nil,
		},
		{
			testCaseName: "case request nil",
			err:          errors.New("request is required"),
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("Case %v : %v", i+1, tc.testCaseName)
		t.Run(name, func(t *testing.T) {
			err := ValidateBroadcastRequest(tc.request)

			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
