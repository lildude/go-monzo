package model

import (
	"database/sql"

	"fmt"

	"github.com/pkg/errors"
)

type Pot struct {
	ID        int64  `sql:"id"`
	PotID     string `json:"id" sql:"pot_id"`
	Name      string `json:"name" sql:"name"`
	AccountID string `sql:"account_id"`
	Balance   int64  `json:"balance" sql:"balance"`
	Currency  string `json:"currency" sql:"currency"`
	Deleted   bool   `json:"deleted"`
}

type Pots struct {
	Array []Pot `json:"pots"`
}

func (pot Pot) Map(rows *sql.Rows) (interface{}, error) {
	err := rows.Scan(
		&pot.ID,
		&pot.PotID,
		&pot.Name,
		&pot.AccountID,
		&pot.Balance,
		&pot.Currency,
		&pot.Deleted,
	)

	return pot, err
}

func (pot Pot) GetID() int64 {
	return pot.ID
}

func (pot Pot) SetID(id int64) interface{} {
	pot.ID = id
	return pot
}

func (pots Pots) ByName(name string) (Pot, error) {
	for _, pot := range pots.Array {
		if pot.Name == name {
			return pot, nil
		}
	}

	return Pot{}, errors.New(fmt.Sprintf("pot not found %s", name))
}
