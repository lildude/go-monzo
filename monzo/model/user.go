package model

import (
	"time"

	"database/sql"
)

type User struct {
	ID           int64     `json:"id" sql:"id"`
	UserID       string    `json:"user_id" sql:"user_id"`
	ClientID     string    `json:"client_id" sql:"client_id"`
	AccessToken  string    `json:"access_token" sql:"access_token"`
	RefreshToken string    `json:"refresh_token" sql:"refresh_token"`
	ExpiresIn    int       `json:"expires_in" sql:"expires_in"`
	ExpiryDate   time.Time `json:"expiry_date" sql:"expiry_date"`
	TokenType    string    `json:"token_type" sql:"token_type"`
}

func (user User) Map(rows *sql.Rows) (interface{}, error) {
	err := rows.Scan(
		&user.ID,
		&user.UserID,
		&user.ClientID,
		&user.AccessToken,
		&user.RefreshToken,
		&user.ExpiryDate,
		&user.TokenType,
		&user.ExpiresIn,
	)

	return user, err
}

func (user User) GetID() int64 {
	return user.ID
}

func (user User) SetID(id int64) interface{} {
	user.ID = id
	return user
}

func (user *User) UpdateExpiry() {
	expiry := time.Second * time.Duration(user.ExpiresIn)
	user.ExpiryDate = time.Now().Add(expiry)
}
