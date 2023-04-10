package entity

type GetBalanceInput struct {
	CallerID             uint64
	PlayerName           string
	Currency             string
	GameID               string
	SessionID            string
	SessionAlternativeID string
	BonusID              string
}

type Balance struct {
	CallerID   uint64
	PlayerName string
	Balance    uint64
	Currency   string

	GameID               string
	SessionID            string
	SessionAlternativeID string
	BonusID              string
	FreeRoundsLeft       int64
}
