package model

type Merchant struct {
	ID         string `json:"id"`
	GroupID    string `json:"group_id"`
	MerchantID string `json:"merchant_id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Logo       string `json:"logo"`
}
