package model

type Pot struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
	Deleted  bool   `json:"deleted"`
}
