package balance

import billingErrors "github.com/Lobanov728/mascot/internal/billing/domain/errors"

func GetDefaultBalanceAmount() uint64 {
	return 0
}

func CalculateNewBalanceAmount(balance uint64, withdraw, deposit int64) (uint64, error) {
	diff := int64(balance) - withdraw + deposit
	if diff < 0 {
		return 0, billingErrors.NewErrorNotEnoughMoney()
	}

	return uint64(diff), nil
}
