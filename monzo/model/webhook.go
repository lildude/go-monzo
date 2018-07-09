package model

import (
	"database/sql"
)

type WebhookBody struct {
	Webhook Webhook `json:"webhook"`
}

type Webhooks struct {
	Array []Webhook `json:"webhooks"`
}

type Webhook struct {
	ID        int64  `sql:"id"`
	AccountID string `json:"account_id" sql:"account_id"`
	WebhookID string `json:"id" sql:"webhook_id"`
	URL       string `json:"url" sql:"url"`
}

func (web Webhook) Map(rows *sql.Rows) (interface{}, error) {
	err := rows.Scan(
		&web.ID,
		&web.AccountID,
		&web.WebhookID,
		&web.URL,
	)

	return web, err
}

func (web Webhook) GetID() int64 {
	return web.ID
}

func (web Webhook) SetID(id int64) interface{} {
	web.ID = id
	return web
}
