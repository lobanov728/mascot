package service

import (
	"context"

	"github.com/Lobanov728/mascot/config"
	"github.com/Lobanov728/mascot/internal/billing/adapters"
	"github.com/Lobanov728/mascot/internal/billing/app"
	"github.com/Lobanov728/mascot/internal/billing/app/command"
	"github.com/Lobanov728/mascot/internal/billing/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context, cfg config.Config) app.Application {
	transactionRepo := adapters.NewTransactionPostgresRepository(cfg.DB.Postgres)
	locker := adapters.NewSimpleLocker()
	logger := logrus.NewEntry(logrus.StandardLogger())
	app := app.Application{
		GetBalanceHandler:         query.NewGetBalanceQueryHandler(transactionRepo, logger),
		ProcessTransactionCommand: command.NewProcessTransactionCommand(transactionRepo, locker, logger),
		ProcessRollbackCommand:    command.NewRollbackCommand(transactionRepo, logger),
	}

	return app
}
