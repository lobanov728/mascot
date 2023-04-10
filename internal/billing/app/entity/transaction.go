package entity

import "time"

type Status string

var (
	RegistredStatus Status = "registred"
	RollbackStatus  Status = "rollback"
	CompletedStatus Status = "completed"
)

type SpinDetail struct {
	BetType string `json:"betType"`
	WinType string `json:"winType"`
}

type Transaction struct {
	CallerID       int64
	PlayerName     string
	Withdraw       int64
	Deposit        int64
	Currency       string
	TransactionRef string
	Status         Status
	CreatedAt      time.Time
	UpdatedAt      time.Time

	GameRoundRef         string
	GameID               string
	Source               string
	Reason               string
	SessionID            string
	SessionAlternativeID string
	SpinDetails          SpinDetail
	BonusID              string
	ChargeFreeRounds     int64
}

type WithdrawAndDepositResponse struct {
	NewBalance     uint64
	TransactionID  string
	FreeRoundsLeft uint64
}
