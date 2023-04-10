package entity

import "github.com/pkg/errors"

var (
	ErrorTransactionNotFound = errors.New("transaction not found")
	ErrorInternalError       = errors.New("internal error")
)

func NewErrorTransactionNotFound() error {
	return ErrorTransactionNotFound
}

func NewInternalError() error {
	return ErrorInternalError
}
