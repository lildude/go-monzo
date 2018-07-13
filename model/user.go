package model

import (
	"time"
)

type User struct {
	UserID       string    `json:"user_id"`
	ClientID     string    `json:"client_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiryDate   time.Time `json:"expiry_date"`
	TokenType    string    `json:"token_type"`
}

func (user *User) UpdateExpiry() {
	expiry := time.Second * time.Duration(user.ExpiresIn)
	user.ExpiryDate = time.Now().Add(expiry)
}
