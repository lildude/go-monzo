package model

import "time"

type Account struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}
