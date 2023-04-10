package command

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/Lobanov728/mascot/internal/billing/app/entity"
// )

// type balanceRepoMock struct {
// 	storage map[string]entity.Balance
// }

// func newBalanceRepoMock(t *testing.T) balanceRepoMock {
// 	t.Helper()

// 	storage := make(map[string]entity.Balance)
// 	storage["abc-1-EUR"] = entity.Balance{
// 		CallerID:   1,
// 		PlayerName: "abc",
// 		Balance:    100,
// 		Currency:   "EUR",
// 	}
// 	m := balanceRepoMock{
// 		storage: storage,
// 	}

// 	return m
// }

// func (m balanceRepoMock) GetBalance(
// 	ctx context.Context,
// 	input entity.GetBalanceInput,
// 	defaultAmount uint64,
// ) (*entity.Balance, error) {
// 	b, ok := m.storage[fmt.Sprintf("%s-%d-%s", input.PlayerName, input.CallerID, input.Currency)]
// 	if !ok {
// 		return nil, nil
// 	}

// 	return &b, nil
// }

// type transactionRepoMock struct {
// 	storage map[string]entity.Balance
// }

// func newTransactionRepoMock(t *testing.T) transactionRepoMock {
// 	t.Helper()

// 	storage := make(map[string]entity.Balance)
// 	storage["abc-1-EUR"] = entity.Balance{
// 		CallerID:   1,
// 		PlayerName: "abc",
// 		Balance:    100,
// 		Currency:   "EUR",
// 	}
// 	m := transactionRepoMock{
// 		storage: storage,
// 	}

// 	return m
// }

// func (m transactionRepoMock) GetTransactionUUIDByTransactionRef(ctx context.Context, transactionRef string) (transactionUUID string, err error) {

// }

// func (m transactionRepoMock) UpdateTransaction(ctx context.Context, uuid string, b *entity.Balance, status entity.Status) error {

// }

// func (m transactionRepoMock) RegisterTransaction(ctx context.Context, t entity.Transaction) (transactionUUID string, err error) {

// }

// func TestProcessTransaction(t *testing.T) {
// 	cmd := NewProcessTransactionCommand()
// }
