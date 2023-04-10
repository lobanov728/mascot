package errors

import "errors"

var (
	ErrorNotEnoughMoney         = errors.New("insufficient funds")
	ErrorNegativeDepositCode    = errors.New("value of deposit ﬁeld is less than zero or incorrect")
	ErrorNegativeWithdrawalCode = errors.New("value of withdraw ﬁeld is less than zero or incorrect")
	ErrorTransactionWasRollback = errors.New("transaction was canceled")
)

func NewErrorNotEnoughMoney() error {
	return ErrorNotEnoughMoney
}

func NewErrorNegativeDeposit() error {
	// Error 3 ErrNegativeDepositCode
	// Value of deposit ﬁeld in the withdrawAndDeposit request is less than zero or incorrect.
	// This value should be greater than zero.
	return ErrorNegativeDepositCode
}

func NewErrorNegativeWithdrawal() error {
	// Error 4 ErrNegativeWithdrawalCode
	// Value of withdraw ﬁeld in the withdrawAndDeposit request is less than zero or incorrect.
	return ErrorNegativeWithdrawalCode
}

func NewErrorTransactionWasRollback() error {
	// In case a rollbackTransaction request has arrived with a transactionRef that is not yet
	// was registered in the service, you need to save the money transaction and mark it as rolled back.
	// If a request for withdrawAndDeposit comes later with the same transactionRef,
	// the service should respond with an error that the transaction could not be committed,
	// because was rolled out
	return ErrorTransactionWasRollback
}
