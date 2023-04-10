package command

import (
	"context"
	"errors"

	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/Lobanov728/mascot/internal/billing/domain/balance"
	"github.com/Lobanov728/mascot/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type RollbackRepository interface {
	GetBalance(ctx context.Context, getBalanceInput entity.GetBalanceInput, defaultAmount uint64) (*entity.Balance, error)
	RegisterRollback(ctx context.Context, rollback entity.Rollback) (uuid string, err error)
	UpdateRollback(
		ctx context.Context,
		uuid string,
		rb entity.Rollback,
		balanceAmount uint64,
	) error
	GetTransactionByTransactionRef(
		ctx context.Context,
		searchInput entity.TransactionSearchInput,
	) (tr *entity.Transaction, err error)
}

type processRollbackCommand struct {
	repo RollbackRepository
}

type ProcessRollbackHandler decorator.CommandHandler[entity.Rollback, *entity.RollbackResponse]

func NewRollbackCommand(r RollbackRepository, logger *logrus.Entry) ProcessRollbackHandler {
	return decorator.ApplyCommandDecorators[entity.Rollback, *entity.RollbackResponse](
		processRollbackCommand{
			repo: r,
		},
		logger,
	)
}

func (c processRollbackCommand) Handle(ctx context.Context, input entity.Rollback) (*entity.RollbackResponse, error) {
	uuid, err := c.repo.RegisterRollback(ctx, entity.Rollback{
		CallerID:             input.CallerID,
		PlayerName:           input.PlayerName,
		Currency:             "",
		RollbackWithdraw:     0,
		RollbackDeposit:      0,
		TransactionRef:       input.TransactionRef,
		GameID:               input.GameID,
		SessionID:            input.SessionID,
		SessionAlternativeID: input.SessionAlternativeID,
		RoundId:              input.RoundId,
	})
	if err != nil {
		return nil, err
	}

	tr, err := c.repo.GetTransactionByTransactionRef(ctx, entity.TransactionSearchInput{
		CallerID:             input.CallerID,
		PlayerName:           input.PlayerName,
		TransactionRef:       input.TransactionRef,
		GameID:               input.GameID,
		SessionID:            input.SessionID,
		SessionAlternativeID: input.SessionAlternativeID,
		RoundId:              input.RoundId,
	})
	if err != nil {
		if errors.Is(err, entity.ErrorTransactionNotFound) {
			return nil, nil
		}
		return nil, err
	}

	b, err := c.repo.GetBalance(ctx, entity.GetBalanceInput{
		CallerID:             input.CallerID,
		PlayerName:           input.PlayerName,
		Currency:             tr.Currency,
		GameID:               input.GameID,
		SessionID:            input.SessionID,
		SessionAlternativeID: input.SessionAlternativeID,
		BonusID:              tr.BonusID,
	}, balance.GetDefaultBalanceAmount())
	if err != nil {
		return nil, err
	}

	balanceAfterRollback, err := balance.CalculateNewBalanceAmount(b.Balance, -tr.Withdraw, -tr.Deposit)
	if err != nil {
		return nil, err
	}

	err = c.repo.UpdateRollback(ctx, uuid, entity.Rollback{
		CallerID:             input.CallerID,
		PlayerName:           input.PlayerName,
		Currency:             tr.Currency,
		RollbackWithdraw:     -tr.Withdraw,
		RollbackDeposit:      -tr.Deposit,
		Status:               entity.CompletedStatus,
		TransactionRef:       input.TransactionRef,
		GameID:               input.GameID,
		SessionID:            input.SessionID,
		SessionAlternativeID: input.SessionAlternativeID,
		RoundId:              input.RoundId,
	}, balanceAfterRollback)

	return nil, err
}
