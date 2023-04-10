package entity

type TransactionSearchInput struct {
	CallerID       uint64
	PlayerName     string
	TransactionRef string

	GameID               string
	SessionID            string
	SessionAlternativeID string
	RoundId              string
}

type Rollback struct {
	CallerID         uint64
	PlayerName       string
	TransactionRef   string
	RollbackWithdraw int64
	RollbackDeposit  int64
	Currency         string
	Status           Status

	GameID               string
	SessionID            string
	SessionAlternativeID string
	RoundId              string
}

type RollbackResponse struct {
}
