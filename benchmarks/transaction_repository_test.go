package benchmarks

import (
	"context"
	"testing"

	"github.com/Lobanov728/mascot/internal/billing/adapters"
	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/google/uuid"
)

func BenchmarkProcessTransactionPlainRepo(b *testing.B) {
	postgres := adapters.Postgres{
		Host:     "localhost",
		Port:     5000,
		Database: "mascot",
		User:     "mascot-user",
		Password: "password",
	}
	repo := adapters.NewTransactionPostgresRepository(postgres)
	ctx := context.Background()

	balance := entity.Balance{
		CallerID:       1,
		PlayerName:     "player_John",
		Balance:        0,
		Currency:       "EUR",
		FreeRoundsLeft: 0,
	}
	for i := 0; i < b.N; i++ {
		repo.RegisterTransaction(ctx, entity.Transaction{
			CallerID:       1,
			PlayerName:     "player_John",
			Withdraw:       100,
			Deposit:        200,
			Currency:       "EUR",
			TransactionRef: uuid.New().String(),
		})

		balance.Balance += 100
	}
}
