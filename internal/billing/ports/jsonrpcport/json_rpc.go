package jsonrpcport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/http/jsonrpc"

	"github.com/Lobanov728/mascot/internal/billing/app"
	"github.com/Lobanov728/mascot/internal/billing/app/entity"
)

func NewJSONRpcServer(app app.Application) *jsonrpc.Server {
	return jsonrpc.NewServer(
		jsonrpc.EndpointCodecMap{
			"getBalance": jsonrpc.EndpointCodec{
				Endpoint: getBalanceEndpoint(app),
				Decode:   decodeGetBalanceRequest,
				Encode:   encodeGetBalanceResponse,
			},
			"withdrawAndDeposit": jsonrpc.EndpointCodec{
				Endpoint: withdrawAndDepositEndpoint(app),
				Decode:   decodeWithdrawAndDepositRequest,
				Encode:   encodeWithdrawAndDepositResponse,
			},
			"rollbackTransaction": jsonrpc.EndpointCodec{
				Endpoint: rollbackTransactionEndpoint(app),
				Decode:   decodeRollbackTransactionRequest,
				Encode:   encodeRollbackTransactionResponse,
			},
		})
}

func getBalanceEndpoint(app app.Application) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*getBalanceRequest) // nolint:errcheck
		b, err := app.GetBalanceHandler.Handle(ctx, entity.GetBalanceInput{
			CallerID:             uint64(req.CallerID),
			PlayerName:           req.PlayerName,
			Currency:             req.Currency,
			GameID:               req.GameID,
			SessionID:            req.SessionID,
			SessionAlternativeID: req.SessionAlternativeID,
			BonusID:              req.BonusID,
		})
		if err != nil {
			return nil, err
		}

		return &getBalanceResponse{
			Balance: b.Balance,
		}, nil
	}
}

func withdrawAndDepositEndpoint(app app.Application) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*withdrawAndDepositRequest) // nolint:errcheck
		resp, err := app.ProcessTransactionCommand.Handle(ctx, entity.Transaction{
			CallerID:             req.CallerID,
			PlayerName:           req.PlayerName,
			Withdraw:             req.Withdraw,
			Deposit:              req.Deposit,
			Currency:             req.Currency,
			TransactionRef:       req.TransactionRef,
			GameRoundRef:         req.GameRoundRef,
			GameID:               req.GameID,
			Source:               req.Source,
			Reason:               req.Reason,
			SessionID:            req.SessionID,
			SessionAlternativeID: req.SessionAlternativeID,
			SpinDetails: entity.SpinDetail{
				BetType: req.SpinDetails.BetType,
				WinType: req.SpinDetails.WinType,
			},
			BonusID:          req.BonusID,
			ChargeFreeRounds: req.ChargeFreeRounds,
		})
		if err != nil {
			return nil, err
		}

		return &withdrawAndDepositResponse{
			NewBalance:     int64(resp.NewBalance),
			TransactionID:  resp.TransactionID,
			FreeRoundsLeft: int64(resp.FreeRoundsLeft),
		}, nil
	}
}

func rollbackTransactionEndpoint(app app.Application) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*rollbackTransactionRequest) // nolint:errcheck
		_, err := app.ProcessRollbackCommand.Handle(ctx, entity.Rollback{
			CallerID:             uint64(req.CallerID),
			PlayerName:           req.PlayerName,
			TransactionRef:       req.TransactionRef,
			GameID:               req.GameID,
			SessionID:            req.SessionID,
			SessionAlternativeID: req.SessionAlternativeID,
			RoundId:              req.RoundId,
		})
		if err != nil {
			return nil, err
		}

		return &rollbackTransactionResponse{}, nil
	}
}
