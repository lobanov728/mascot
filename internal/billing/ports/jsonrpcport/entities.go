package jsonrpcport

type getBalanceResponse struct {
	Balance        uint64 `json:"balance"`
	FreeRoundsLeft uint64 `json:"freeroundsLeft,omitempty"`
}

type getBalanceRequest struct {
	PlayerName           string `json:"playerName"`
	Currency             string `json:"currency"`
	CallerID             int64  `json:"callerId"`
	GameID               string `json:"gameId,omitempty"`
	SessionID            string `json:"sessionId,omitempty"`
	SessionAlternativeID string `json:"sessionAlternativeID,omitempty"`
	BonusID              string `json:"bonusID,omitempty"`
}

type withdrawAndDepositResponse struct {
	NewBalance     int64  `json:"newBalance"`
	TransactionID  string `json:"transactionId"`
	FreeRoundsLeft int64  `json:"freeroundsLeft"`
}

type withdrawAndDepositRequest struct {
	CallerID             int64       `json:"callerId"`
	PlayerName           string      `json:"playerName"`
	Withdraw             int64       `json:"withdraw"`
	Deposit              int64       `json:"deposit"`
	Currency             string      `json:"currency"`
	TransactionRef       string      `json:"transactionRef"`
	GameRoundRef         string      `json:"gameRoundRef"`
	GameID               string      `json:"gameID"`
	Source               string      `json:"source"`
	Reason               string      `json:"reasin"`
	SessionID            string      `json:"sessionId"`
	SessionAlternativeID string      `json:"sessionAlternativeId"`
	SpinDetails          *spinDetail `json:"spinDetails,omitempty"`
	BonusID              string      `json:"bonusID"`
	ChargeFreeRounds     int64       `json:"chargeFreeRounds"`
}

type spinDetail struct {
	BetType string `json:"betType,omitempty"`
	WinType string `json:"winType,omitempty"`
}

type rollbackTransactionResponse struct {
}

type rollbackTransactionRequest struct {
	PlayerName           string `json:"playerName"`
	CallerID             int64  `json:"callerId"`
	TransactionRef       string `json:"transactionRef"`
	GameID               string `json:"gameId,omitempty"`
	SessionID            string `json:"sessionId,omitempty"`
	SessionAlternativeID string `json:"sessionAlternativeID,omitempty"`
	RoundId              string `json:"roundId,omitempty"`
}
