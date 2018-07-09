package model

import "time"

type Transaction struct {
	Type string `json:"type"`
	Data struct {
		TransactionID string    `json:"id"`
		AccountID     string    `json:"account_id"`
		Description   string    `json:"description"`
		Category      string    `json:"category"`
		Amount        int64     `json:"amount"`
		Currency      string    `json:"currency"`
		Created       time.Time `json:"created"`
		IsLoad        bool      `json:"is_load"`
		Merchant      struct {
			ID       string `json:"id"`
			GroupID  string `json:"group_id"`
			Name     string `json:"name"`
			Category string `json:"category"`
			Logo     string `json:"logo"`
		} `json:"merchant"`
	} `json:"data"`
}
