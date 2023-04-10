package app

import (
	"github.com/Lobanov728/mascot/internal/billing/app/command"
	"github.com/Lobanov728/mascot/internal/billing/app/query"
)

type Application struct {
	GetBalanceHandler         query.GetBalanceHandler
	ProcessRollbackCommand    command.ProcessRollbackHandler
	ProcessTransactionCommand command.ProcessTransactionHangler
}
