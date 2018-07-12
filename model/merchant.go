package model

type Merchant struct {
	ID         int64  `json:"id"`
	GroupID    string `json:"group_id"`
	MerchantID string `json:"merchant_id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Logo       string `json:"logo"`
}
