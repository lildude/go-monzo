package model

type Webhook struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	URL       string `json:"url"`
}
