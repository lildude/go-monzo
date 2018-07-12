package model

type Balance struct {
	Balance      int64  `json:"balance"`
	TotalBalance int64  `json:"total_balance"`
	SpendToday   int64  `json:"spend_today"`
	Currency     string `json:"currency"`
}
