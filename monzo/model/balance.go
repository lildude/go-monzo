package model

type Balance struct {
	Balance      int64  `json:"balance"`
	TotalBalance int64  `json:"total_balance"`
	Currency     string `json:"currency"`
}
