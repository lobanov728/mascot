package adapters

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/Lobanov728/mascot/internal/billing/app/entity"
)

type SpinDetail struct {
	BetType sql.NullString
	WinType sql.NullString
}

type Transaction struct {
	TransactionUUID      string
	TransactionRef       string
	CreatedAt            time.Time
	UpdatedAt            sql.NullTime
	CallerID             int64
	PlayerUUID           string
	Withdraw             uint64
	Deposit              uint64
	Status               string
	Balance              uint64
	Currency             string
	GameID               sql.NullString
	GameRoundRef         sql.NullString
	SessionID            sql.NullString
	SessionAlternativeID sql.NullString
	Source               sql.NullString
	Reason               sql.NullString
	SpinDetails          SpinDetail
	BonusID              sql.NullString
	FreeRoundsLeft       sql.NullInt64
	ChargeFreeRounds     sql.NullInt64
}

type Balance struct {
	CallerID             int64
	PlayerUUID           string
	Balance              uint64
	Currency             string
	GameID               sql.NullString
	SessionID            sql.NullString
	SessionAlternativeID sql.NullString
	BonusID              sql.NullString
	FreeRoundsLeft       sql.NullInt64
}

type TransactionPostgresRepository struct {
	conn *sql.DB
}

func NewTransactionPostgresRepository(dbConfig Postgres) *TransactionPostgresRepository {
	conn, err := sql.Open("postgres", dbConfig.URL())

	if err != nil {
		panic(fmt.Sprintf("fail to open postgress connection: %s", err))
	}

	if err := conn.Ping(); err != nil {
		panic(fmt.Sprintf("fail to ping postgress connection: %s", err))
	}

	return &TransactionPostgresRepository{
		conn: conn,
	}
}

func (r *TransactionPostgresRepository) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	tx, err = r.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *TransactionPostgresRepository) Commit(tx *sql.Tx) (err error) {
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

func (r *TransactionPostgresRepository) Rolback(tx *sql.Tx) (err error) {
	if tx != nil {
		return tx.Rollback()
	}

	return nil
}

func (r *TransactionPostgresRepository) GetBalance(ctx context.Context, getBalanceInput entity.GetBalanceInput, defaultAmount uint64) (*entity.Balance, error) {
	query := `SELECT 
	caller_id,
	player_uuid,
	balance,
	currency,
	game_id,
	session_id,
	session_alternative_id,
	bonus_id,
	free_rounds_left
	FROM transaction t
	WHERE t.player_uuid = $1 AND t.caller_id = $2 AND t.currency = $3 AND t.status=$4 ORDER BY t.created_at DESC LIMIT 1`
	row := r.conn.QueryRow(query, getBalanceInput.PlayerName, getBalanceInput.CallerID, getBalanceInput.Currency, entity.CompletedStatus)

	var balance Balance

	err := row.Scan(
		&balance.CallerID,
		&balance.PlayerUUID,
		&balance.Balance,
		&balance.Currency,
		&balance.GameID,
		&balance.SessionID,
		&balance.SessionAlternativeID,
		&balance.BonusID,
		&balance.FreeRoundsLeft,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.Balance{
				CallerID:             uint64(balance.CallerID),
				PlayerName:           balance.PlayerUUID,
				Balance:              balance.Balance,
				Currency:             balance.Currency,
				GameID:               balance.GameID.String,
				SessionID:            balance.SessionID.String,
				SessionAlternativeID: balance.SessionAlternativeID.String,
				BonusID:              balance.BonusID.String,
				FreeRoundsLeft:       balance.FreeRoundsLeft.Int64,
			}, nil
		}

		return nil, err
	}

	return &entity.Balance{
		CallerID:             uint64(balance.CallerID),
		PlayerName:           balance.PlayerUUID,
		Balance:              balance.Balance,
		Currency:             balance.Currency,
		GameID:               balance.GameID.String,
		SessionID:            balance.SessionID.String,
		SessionAlternativeID: balance.SessionAlternativeID.String,
		BonusID:              balance.BonusID.String,
		FreeRoundsLeft:       balance.FreeRoundsLeft.Int64,
	}, nil
}

func (r *TransactionPostgresRepository) HasUnproccessedRollbacks(ctx context.Context, transactionRef string) (bool, error) {
	query := `SELECT count(*) FROM transaction WHERE transaction_ref = $1 AND status = $2;`
	row := r.conn.QueryRowContext(ctx, query, transactionRef, entity.RollbackStatus)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *TransactionPostgresRepository) RegisterRollback(ctx context.Context, rollback entity.Rollback) (res string, err error) {
	newUUID := uuid.NewString()
	query := `INSERT INTO transaction
	(
	transaction_uuid,
	transaction_ref,
	caller_id,
	player_uuid,
	withdraw,
	deposit,
	status,
	balance,
	currency,
	game_id,
	session_id,
	session_alternative_id) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
	);
	`
	_, err = r.conn.ExecContext(
		ctx,
		query,
		newUUID,
		rollback.TransactionRef,
		rollback.CallerID,
		rollback.PlayerName,
		0,
		0,
		entity.RollbackStatus,
		0,
		rollback.Currency,
		rollback.GameID,
		rollback.SessionID,
		rollback.SessionAlternativeID,
	)

	if err != nil {
		return "", err
	}

	return newUUID, nil
}

func (r *TransactionPostgresRepository) UpdateRollback(
	ctx context.Context,
	uuid string,
	rollback entity.Rollback,
	balanceAmount uint64,
) error {
	query := `UPDATE transaction
	SET withdraw=$1, deposit=$2, balance=$3, currency=$4, status=$5
	WHERE transaction_uuid=$6;
	`
	_, err := r.conn.ExecContext(
		ctx,
		query,
		rollback.RollbackWithdraw,
		rollback.RollbackDeposit,
		balanceAmount,
		rollback.Currency,
		rollback.Status,
		uuid,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionPostgresRepository) GetTransactionUUIDByTransactionRef(ctx context.Context, transactionRef string) (transactionUUID string, err error) {
	query := `SELECT transaction_uuid FROM transaction t WHERE t.transaction_ref = $1 and t.status=$2 LIMIT 1`
	rows := r.conn.QueryRow(query, transactionRef, entity.CompletedStatus)

	err = rows.Scan(&transactionUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", err
	}

	return transactionUUID, nil
}

func (r *TransactionPostgresRepository) GetTransactionByTransactionRef(ctx context.Context, searchInput entity.TransactionSearchInput) (tr *entity.Transaction, err error) {
	query := `SELECT 
	transaction_uuid,
	transaction_ref,
	caller_id,
	player_uuid,
	withdraw,
	deposit,
	balance,
	currency,
	game_id,
	game_round_ref,
	source,
	reason,
	session_id,
	session_alternative_id,
	bonus_id,
	charge_free_rounds,
	free_rounds_left
	FROM transaction t
	WHERE t.caller_id = $1 AND t.player_uuid = $2 AND t.transaction_ref = $3 and (deposit > 0 or withdraw > 0)
	`
	query += ` ORDER BY t.created_at DESC LIMIT 1 `
	row := r.conn.QueryRow(
		query,
		searchInput.CallerID,
		searchInput.PlayerName,
		searchInput.TransactionRef,
	)
	var transaction Transaction
	err = row.Scan(
		&transaction.TransactionUUID,
		&transaction.TransactionRef,
		&transaction.CallerID,
		&transaction.PlayerUUID,
		&transaction.Withdraw,
		&transaction.Deposit,
		&transaction.Balance,
		&transaction.Currency,
		&transaction.GameID,
		&transaction.GameRoundRef,
		&transaction.Source,
		&transaction.Reason,
		&transaction.SessionID,
		&transaction.SessionAlternativeID,
		&transaction.BonusID,
		&transaction.ChargeFreeRounds,
		&transaction.FreeRoundsLeft,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(entity.NewErrorTransactionNotFound(), err.Error())
		}

		return nil, err
	}

	return &entity.Transaction{
		CallerID:             transaction.CallerID,
		PlayerName:           transaction.PlayerUUID,
		Withdraw:             int64(transaction.Withdraw),
		Deposit:              int64(transaction.Deposit),
		Currency:             transaction.Currency,
		TransactionRef:       transaction.TransactionRef,
		CreatedAt:            time.Time{},
		UpdatedAt:            time.Time{},
		GameRoundRef:         transaction.GameRoundRef.String,
		GameID:               transaction.GameID.String,
		Source:               transaction.Source.String,
		Reason:               transaction.Reason.String,
		SessionID:            transaction.SessionID.String,
		SessionAlternativeID: transaction.SessionAlternativeID.String,
		BonusID:              transaction.BonusID.String,
		ChargeFreeRounds:     transaction.ChargeFreeRounds.Int64,
	}, nil
}

func (r TransactionPostgresRepository) RegisterTransaction(ctx context.Context, t entity.Transaction) (transactionUUID string, err error) {
	uuid := uuid.NewString()
	query := `INSERT INTO transaction
	(
	transaction_uuid,
	transaction_ref,
	caller_id,
	player_uuid,
	withdraw,
	deposit,
	balance,
	status,
	currency,
	game_id,
	game_round_ref,
	source,
	reason,
	session_id,
	session_alternative_id,
	spin_details,
	bonus_id,
	charge_free_rounds) VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
	);
	`
	spinDetails, err := json.Marshal(t.SpinDetails)
	if err != nil {
		return "", err
	}

	_, err = r.conn.ExecContext(
		ctx,
		query,
		uuid,
		t.TransactionRef,
		t.CallerID,
		t.PlayerName,
		t.Withdraw,
		t.Deposit,
		0,
		t.Status,
		t.Currency,
		t.GameID,
		t.GameRoundRef,
		t.Source,
		t.Reason,
		t.SessionID,
		t.SessionAlternativeID,
		string(spinDetails),
		t.BonusID,
		t.ChargeFreeRounds,
	)

	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (r TransactionPostgresRepository) UpdateTransaction(ctx context.Context, uuid string, b *entity.Balance, status entity.Status) error {
	query := `UPDATE transaction SET  balance=$1, free_rounds_left=$2, status=$3 WHERE transaction_uuid = $4;`
	_, err := r.conn.ExecContext(
		ctx,
		query,
		b.Balance,
		b.FreeRoundsLeft,
		status,
		uuid,
	)

	return err
}
