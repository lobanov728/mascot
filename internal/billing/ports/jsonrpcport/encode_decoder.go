package jsonrpcport

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func decodeGetBalanceRequest(ctx context.Context, msg json.RawMessage) (interface{}, error) {
	var req getBalanceRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeGetBalanceResponse(ctx context.Context, result interface{}) (json.RawMessage, error) {
	res, ok := result.(*getBalanceResponse)
	if !ok {
		return nil, errors.New("wrong response type")
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func decodeWithdrawAndDepositRequest(ctx context.Context, msg json.RawMessage) (interface{}, error) {
	var req withdrawAndDepositRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeWithdrawAndDepositResponse(ctx context.Context, result interface{}) (json.RawMessage, error) {
	res, ok := result.(*withdrawAndDepositResponse)
	if !ok {
		return nil, errors.New("wrong response type")
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func decodeRollbackTransactionRequest(ctx context.Context, msg json.RawMessage) (interface{}, error) {
	var req rollbackTransactionRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeRollbackTransactionResponse(ctx context.Context, result interface{}) (json.RawMessage, error) {
	res, ok := result.(*rollbackTransactionResponse)
	if !ok {
		return nil, errors.New("wrong response type")
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}
