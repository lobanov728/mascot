package command

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/Lobanov728/mascot/internal/billing/domain/balance"
	billingErrors "github.com/Lobanov728/mascot/internal/billing/domain/errors"
	domainTransaction "github.com/Lobanov728/mascot/internal/billing/domain/transaction"
	"github.com/Lobanov728/mascot/internal/common/decorator"
)

type TransactionRepository interface {
	BeginTx(ctx context.Context) (tx *sql.Tx, err error)
	Commit(tx *sql.Tx) (err error)
	Rolback(tx *sql.Tx) (err error)
	GetBalance(ctx context.Context, getBalanceInput entity.GetBalanceInput, defaultAmount uint64) (*entity.Balance, error)
	HasUnproccessedRollbacks(ctx context.Context, transactionRef string) (bool, error)
	GetTransactionUUIDByTransactionRef(ctx context.Context, transactionRef string) (transactionUUID string, err error)
	UpdateTransaction(ctx context.Context, uuid string, b *entity.Balance, status entity.Status) error
	RegisterTransaction(ctx context.Context, t entity.Transaction) (transactionUUID string, err error)
}

type processTransactionCommand struct {
	locker Locker
	repo   TransactionRepository
	mu     *sync.Mutex
}

type ProcessTransactionHangler decorator.CommandHandler[entity.Transaction, *entity.WithdrawAndDepositResponse]

func NewProcessTransactionCommand(
	r TransactionRepository,
	l Locker,
	logger *logrus.Entry,
) ProcessTransactionHangler {

	return decorator.ApplyCommandDecorators[entity.Transaction, *entity.WithdrawAndDepositResponse](
		processTransactionCommand{
			repo:   r,
			locker: l,
			mu:     &sync.Mutex{},
		},
		logger,
	)
}

func (c processTransactionCommand) Handle(
	ctx context.Context,
	t entity.Transaction,
) (*entity.WithdrawAndDepositResponse, error) {
	// lock by transaction ref protect from simultaneous inputs with the same transaction
	if err := domainTransaction.Validate(t.Deposit, t.Withdraw); err != nil {
		return nil, err
	}

	// l, err := c.locker.AcuireLock(ctx, t.PlayerName, time.Millisecond*50)
	// if l == nil && err == nil {
	// 	return &entity.WithdrawAndDepositResponse{
	// 		NewBalance:     0,
	// 		TransactionID:  "123",
	// 		FreeRoundsLeft: uint64(0),
	// 	}, nil
	// }
	// defer c.locker.Release(ctx, l)

	// Error 2 ErrIllegalCurrencyCode
	// Game Provider sends a request with a currency code that does not support by the server on the operator’s
	// side. Один игрок – одна валюта

	// register transaction
	t.Status = entity.RegistredStatus
	transactionUUID, err := c.repo.RegisterTransaction(ctx, t)
	if err != nil {
		// need to log error
		return nil, err
	}

	// check on duplicates
	duplicateUUID, err := c.repo.GetTransactionUUIDByTransactionRef(ctx, t.TransactionRef)
	if err != nil {
		// if got error belive that duplicate doesn't exist
		duplicateUUID = ""
		// need to log error
	}

	if duplicateUUID != "" {
		dB, err := c.repo.GetBalance(ctx, entity.GetBalanceInput{
			CallerID:             uint64(t.CallerID),
			PlayerName:           t.PlayerName,
			Currency:             t.Currency,
			GameID:               t.GameID,
			SessionID:            t.SessionID,
			SessionAlternativeID: t.SessionAlternativeID,
			BonusID:              t.BonusID,
		}, balance.GetDefaultBalanceAmount())
		if err != nil {
			// need to log error
			return &entity.WithdrawAndDepositResponse{
				NewBalance:     0,
				TransactionID:  duplicateUUID,
				FreeRoundsLeft: uint64(0),
			}, nil
		}
		return &entity.WithdrawAndDepositResponse{
			NewBalance:     dB.Balance,
			TransactionID:  duplicateUUID,
			FreeRoundsLeft: uint64(dB.FreeRoundsLeft),
		}, nil
	}

	// check on rollbacks
	hasRollbacks, err := c.repo.HasUnproccessedRollbacks(ctx, t.TransactionRef)
	if err != nil {
		return nil, err
	}

	if hasRollbacks {
		return nil, billingErrors.NewErrorTransactionWasRollback()
	}

	// if no duplicate no rollbacks save transaction make changes with balance

	tx, err := c.repo.BeginTx(ctx)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	bal, err := c.repo.GetBalance(ctx, entity.GetBalanceInput{
		CallerID:             uint64(t.CallerID),
		PlayerName:           t.PlayerName,
		Currency:             t.Currency,
		GameID:               t.GameID,
		SessionID:            t.SessionID,
		SessionAlternativeID: t.SessionAlternativeID,
		BonusID:              t.BonusID,
	}, balance.GetDefaultBalanceAmount())
	if err != nil {
		fmt.Println(err)
		// need to log error
		return &entity.WithdrawAndDepositResponse{
			NewBalance:     0,
			TransactionID:  duplicateUUID,
			FreeRoundsLeft: uint64(0),
		}, nil
	}

	newBalanceAmount, err := balance.CalculateNewBalanceAmount(bal.Balance, t.Withdraw, t.Deposit)
	if err != nil && errors.Is(err, billingErrors.ErrorNotEnoughMoney) {
		return nil, err
	}

	bal.Balance = newBalanceAmount
	err = c.repo.UpdateTransaction(ctx, transactionUUID, bal, entity.CompletedStatus)
	if err != nil {
		fmt.Println(err)
		// log error
		return &entity.WithdrawAndDepositResponse{
			NewBalance:     newBalanceAmount,
			TransactionID:  duplicateUUID,
			FreeRoundsLeft: uint64(0),
		}, nil
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		// log error
		return nil, err
	}
	return &entity.WithdrawAndDepositResponse{
		NewBalance:     newBalanceAmount,
		TransactionID:  transactionUUID,
		FreeRoundsLeft: uint64(bal.FreeRoundsLeft),
	}, nil
}
