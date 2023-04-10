package benchmarks

import (
	"context"
	"sync"
	"testing"

	"github.com/Lobanov728/mascot/internal/billing/adapters"
	"github.com/Lobanov728/mascot/internal/billing/app/command"
	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestProcessTransactionInConcurentWay(t *testing.T) {
	t.Parallel()

	postgres := adapters.Postgres{
		Host:     "localhost",
		Port:     5000,
		Database: "mascot",
		User:     "mascot-user",
		Password: "password",
	}

	balanceRepo := adapters.NewTransactionPostgresRepository(postgres)
	locker := adapters.NewSimpleLocker()

	cmd := command.NewProcessTransactionCommand(balanceRepo, locker, nil)

	ctx := context.Background()

	tr := entity.Transaction{
		PlayerName: "player1",
		CallerID:   1,
		Withdraw:   100,
		Deposit:    200,
		Currency:   "EUR",
	}
	respCh1 := make(chan *entity.WithdrawAndDepositResponse)
	respCh2 := make(chan *entity.WithdrawAndDepositResponse)

	go func() {
		for i := 0; i < 10; i++ {
			copyTr := tr
			copyTr.TransactionRef = uuid.NewString()
			resp, err := cmd.Handle(ctx, copyTr)
			require.NoError(t, err)
			respCh1 <- resp
		}
		close(respCh1)
	}()
	go func() {
		for i := 0; i < 10; i++ {
			copyTr := tr
			copyTr.TransactionRef = uuid.NewString()
			resp, err := cmd.Handle(ctx, copyTr)
			require.NoError(t, err)
			respCh2 <- resp
		}
		close(respCh2)
	}()

	balances := make(map[uint64]int, 20)
	for resp := range merge(respCh1, respCh2) {
		balances[resp.NewBalance]++
	}

	for i := range balances {
		require.Equal(t, 1, balances[i])
	}
}

func merge(chans ...chan *entity.WithdrawAndDepositResponse) chan *entity.WithdrawAndDepositResponse {
	respCh := make(chan *entity.WithdrawAndDepositResponse)
	var wg sync.WaitGroup

	for _, c := range chans {
		wg.Add(1)
		c := c
		go func() {
			defer wg.Done()
			for resp := range c {
				respCh <- resp
			}
		}()
	}

	go func() {
		wg.Wait()
		close(respCh)
	}()

	return respCh
}
