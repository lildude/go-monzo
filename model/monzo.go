package model

import (
	"errors"
	"fmt"
)

type Monzo struct {
	Accounts []Account `json:"accounts"`
	Pots     []Pot     `json:"pots"`
	Webhooks []Webhook `json:"webhooks"`
	Webhook  Webhook   `json:"webhook"`
}

func (monzo Monzo) ByName(name string) (Pot, error) {
	for _, pot := range monzo.Pots {
		if pot.Name == name {
			return pot, nil
		}
	}

	return Pot{}, errors.New(fmt.Sprintf("pot not found %s", name))
}
