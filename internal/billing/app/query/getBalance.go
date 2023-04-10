package query

import (
	"context"

	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/Lobanov728/mascot/internal/billing/domain/balance"
	"github.com/Lobanov728/mascot/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type GetBalanceRepository interface {
	GetBalance(ctx context.Context, getBalanceInput entity.GetBalanceInput, defaultAmount uint64) (*entity.Balance, error)
}

type getBalanceHandler struct {
	BalanceRepo GetBalanceRepository
}

type GetBalanceHandler decorator.QueryHandler[entity.GetBalanceInput, *entity.Balance]

func NewGetBalanceQueryHandler(repo GetBalanceRepository, logger *logrus.Entry) GetBalanceHandler {
	return decorator.ApplyQueryDecorators[entity.GetBalanceInput, *entity.Balance](
		getBalanceHandler{
			BalanceRepo: repo,
		},
		logger,
	)
}

func (q getBalanceHandler) Handle(ctx context.Context, input entity.GetBalanceInput) (*entity.Balance, error) {
	b, err := q.BalanceRepo.GetBalance(ctx, input, balance.GetDefaultBalanceAmount())
	if err != nil {
		return nil, err
	}

	return b, nil
}
