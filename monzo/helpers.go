package monzo

import (
	"github.com/gurparit/go-monzo/monzo/model"
	"github.com/gurparit/callisto/sqlio"
)

func GetUser(userID string) (model.User, error) {
	sqlioUser := sqlio.New(model.User{})
	if err := sqlioUser.SelectWhere(sqlio.Values{"user_id": userID}); err != nil {
		return model.User{}, err
	}

	return sqlioUser.Get().(model.User), nil
}

func GetAccountByUser(userID string) (model.Account, error) {
	sqlioAccount := sqlio.New(model.Account{})
	if err := sqlioAccount.SelectWhere(sqlio.Values{"user_id": userID}); err != nil {
		return model.Account{}, err
	}

	return sqlioAccount.Get().(model.Account), nil
}

func GetAccountByID(accountID string) (model.Account, error) {
	sqlioAccount := sqlio.New(model.Account{})
	if err := sqlioAccount.SelectWhere(sqlio.Values{"account_id": accountID}); err != nil {
		return model.Account{}, err
	}

	return sqlioAccount.Get().(model.Account), nil
}

func DeleteWebhookByID(webhookID string) error {
	sqlioWebhook := sqlio.New(model.Webhook{})
	if err := sqlioWebhook.SelectWhere(sqlio.Values{"webhook_id": webhookID}); err != nil {
		return err
	} else {
		sqlioWebhook.Delete()
	}

	return nil
}

func SaveOrUpdateUser(user model.User) {
	sql := sqlio.New(model.User{})
	if err := sql.SelectWhere(sqlio.Values{"user_id": user.UserID}); err != nil {
		sql = sqlio.New(user)
		sql.Save()
	} else {
		sql.Update(user)
	}
}

func SaveOrUpdateAccount(account model.Account) {
	sql := sqlio.New(model.Account{})
	if err := sql.SelectWhere(sqlio.Values{"account_id": account.AccountID}); err != nil {
		sql = sqlio.New(account)
		sql.Save()
	} else {
		sql.Update(account)
	}
}

func SaveWebhook(webhook model.Webhook) {
	sqlioWebhook := sqlio.New(webhook)
	sqlioWebhook.Save()
}
