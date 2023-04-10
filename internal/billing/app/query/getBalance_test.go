package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/Lobanov728/mascot/internal/billing/app/entity"
	"github.com/stretchr/testify/require"
)

type balanceRepoMock struct {
	storage map[string]entity.Balance
}

func newBalanceRepoMock(t *testing.T) balanceRepoMock {
	t.Helper()

	storage := make(map[string]entity.Balance)
	storage["abc-1-EUR"] = entity.Balance{
		CallerID:   1,
		PlayerName: "abc",
		Balance:    100,
		Currency:   "EUR",
	}
	storage["abcd-1-EUR"] = entity.Balance{
		CallerID:   1,
		PlayerName: "abcd",
		Balance:    200,
		Currency:   "EUR",
	}
	m := balanceRepoMock{
		storage: storage,
	}

	return m
}

func (m balanceRepoMock) GetBalance(
	ctx context.Context,
	input entity.GetBalanceInput,
	defaultAmount uint64,
) (*entity.Balance, error) {
	b, ok := m.storage[fmt.Sprintf("%s-%d-%s", input.PlayerName, input.CallerID, input.Currency)]
	if !ok {
		return nil, nil
	}

	return &b, nil
}

func TestGetBalance(t *testing.T) {
	t.Parallel()

	query := NewGetBalanceQueryHandler(newBalanceRepoMock(t), nil)
	ctx := context.Background()
	b, err := query.Handle(ctx, entity.GetBalanceInput{
		CallerID:   1,
		PlayerName: "abc",
		Currency:   "EUR",
	})

	require.NoError(t, err)
	require.Equal(t, uint64(100), b.Balance)
}
