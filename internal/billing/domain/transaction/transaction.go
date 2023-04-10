package transaction

import (
	"time"

	billingErrors "github.com/Lobanov728/mascot/internal/billing/domain/errors"
)

type Transaction struct {
	CallerID       int64
	PlayerName     string
	Currency       string
	TransactionRef string
	CreatedAt      time.Time
}

func Validate(deposit int64, withdraw int64) error {
	if deposit < 0 {
		return billingErrors.NewErrorNegativeDeposit()
	}

	if withdraw < 0 {
		return billingErrors.NewErrorNegativeWithdrawal()
	}

	return nil
}
