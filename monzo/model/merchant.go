package model

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Merchant struct {
	ID         int64  `json:"id"`
	GroupID    string `json:"group_id"`
	MerchantID string `json:"merchant_id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	Logo       string `json:"logo"`
}

func (m Merchant) Map(rows *sql.Rows) (interface{}, error) {
	return nil, errors.New("not implemented error")
}

func (m Merchant) GetID() int64 {
	return m.ID
}

func (m Merchant) SetID(id int64) interface{} {
	m.ID = id
	return m
}
