package model

import (
	"database/sql"
)

type Account struct {
	ID          int64  `json:"-" sql:"id"`
	UserID      string `json:"user_id" sql:"user_id"`
	AccountID   string `json:"id" sql:"account_id"`
	Description string `json:"description" sql:"description"`
}

type Accounts struct {
	Array []Account `json:"accounts"`
}

func (acc Account) Map(rows *sql.Rows) (interface{}, error) {
	err := rows.Scan(
		&acc.ID,
		&acc.UserID,
		&acc.AccountID,
		&acc.Description,
	)

	return acc, err
}

func (acc Account) GetID() int64 {
	return acc.ID
}

func (acc Account) SetID(id int64) interface{} {
	acc.ID = id
	return acc
}
