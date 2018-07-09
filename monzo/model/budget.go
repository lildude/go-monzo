package model

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Budget struct {
	ID        int64  `json:"id" sql:"id"`
	AccountID string `json:"account_id" sql:"account_id"`
	GroupID   string `json:"group_id" sql:"group_id"`
	Merchant  string `json:"merchant_name" sql:"merchant_name"`
	PotID     string `json:"pot_id" sql:"pot_id"`
	PotName   string `json:"pot_name" sql:"pot_name"`
	Currency  string `json:"currency" sql:"currency"`
}

func (b Budget) Map(rows *sql.Rows) (interface{}, error) {
	return nil, errors.New("not implemented error")
}

func (b Budget) GetID() int64 {
	return b.ID
}

func (b Budget) SetID(id int64) interface{} {
	b.ID = id
	return b
}
