package monzo

import (
	"net/http"

	"strconv"

	"os"

	"fmt"

	"strings"

	"github.com/gurparit/go-monzo/model"
	"github.com/gurparit/go-common/httpc"
	"github.com/gurparit/go-common/logio"
	"github.com/gurparit/go-common/uuid"
)

const (
	monzoClientID     = "MONZO_CLIENT_ID"
	monzoClientSecret = "MONZO_CLIENT_SECRET"
	monzoRedirectURI  = "MONZO_REDIRECT_URI"
	monzoWebhookURI   = "MONZO_WEBHOOK_URI"

	LoginURL          = "MONZO_URL_LOGIN"
	WhoAmIURL         = "MONZO_URL_WHOAMI"
	Oauth2URL         = "MONZO_URL_OAUTH2"
	AccountsURL       = "MONZO_URL_ACCOUNTS"
	BalanceURL        = "MONZO_URL_BALANCE"
	PotsURL           = "MONZO_URL_POTS"
	DepositURL        = "MONZO_URL_POTS_DEPOSIT"
	WithdrawURL       = "MONZO_URL_POTS_WITHDRAW"
	WebhookGetURL     = "MONZO_URL_WEBHOOK"
	WebhookCreateURL  = "MONZO_URL_WEBHOOK_CREATE"
	WebhookDeleteURL  = "MONZO_URL_WEBHOOK_DELETE"
	FeedItemCreateURL = "MONZO_URL_FEED"
)

var (
	ClientID     = os.Getenv(monzoClientID)
	ClientSecret = os.Getenv(monzoClientSecret)
	RedirectURI  = os.Getenv(monzoRedirectURI)
	WebhookURI   = os.Getenv(monzoWebhookURI)
)

var urlMap = map[string]string{
	LoginURL:          "https://auth.monzo.com/?client_id=%s&redirect_uri=%s&response_type=code&state=%s",
	WhoAmIURL:         "https://api.monzo.com/ping/whoami",
	Oauth2URL:         "https://api.monzo.com/oauth2/token",
	AccountsURL:       "https://api.monzo.com/accounts?account_type=uk_retail",
	BalanceURL:        "https://api.monzo.com/balance?account_id=%s",
	PotsURL:           "https://api.monzo.com/pots",
	DepositURL:        "https://api.monzo.com/pots/%s/deposit",
	WithdrawURL:       "https://api.monzo.com/pots/%s/withdraw",
	WebhookGetURL:     "https://api.monzo.com/webhooks?account_id=%s",
	WebhookCreateURL:  "https://api.monzo.com/webhooks",
	WebhookDeleteURL:  "https://api.monzo.com/webhooks/%s",
	FeedItemCreateURL: "https://api.monzo.com/feed",
}

type Monzo struct {
	tokenType string
	accessToken string
}

func GetURL(urlname string, params ...interface{}) string {
	baseURL := urlMap[urlname]
	if strings.Contains(baseURL, "%s") {
		return fmt.Sprintf(baseURL, params...)
	}

	return baseURL
}

func SetURL(urlname string, newURL string) {
	urlMap[urlname] = newURL
}

func New(tokenType string, accessToken string) Monzo {
	return Monzo{
		tokenType: tokenType,
		accessToken: accessToken,
	}
}

func Login(state string) string {
	return GetURL(LoginURL, ClientID, RedirectURI, state)
}

func Callback(code string) (model.User, error) {
	headers := make(httpc.Headers)
	headers.FormURLEncoded()

	data := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     ClientID,
		"client_secret": ClientSecret,
		"redirect_uri":  RedirectURI,
		"code":          code,
	}

	targetURL := GetURL(Oauth2URL)
	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodPost,
		Headers:   headers,
		Form:      data,
	}

	var user model.User
	if err := request.JSON(&user); err != nil {
		return model.User{}, err
	}

	user.UpdateExpiry()

	return user, nil
}

func Refresh(refreshToken string) (model.User, error) {
	headers := make(httpc.Headers)
	headers.FormURLEncoded()

	data := make(map[string]string)
	data["grant_type"] = "refresh_token"
	data["client_id"] = ClientID
	data["client_secret"] = ClientSecret
	data["refresh_token"] = refreshToken

	request := httpc.HTTP{
		TargetURL: GetURL(Oauth2URL),
		Method:    http.MethodPost,
		Headers:   headers,
		Form:      data,
	}

	var user model.User
	if err := request.JSON(&user); err != nil {
		return model.User{}, err
	} else {
		user.UpdateExpiry()
	}

	return user, nil
}

func (m Monzo) WhoAmI() (model.WhoAmI, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)
	headers.FormURLEncoded()

	targetURL := GetURL(WhoAmIURL)

	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodPut,
		Headers:   headers,
		Form:      nil,
	}

	var whoami model.WhoAmI
	if err := request.JSON(&whoami); err != nil {
		return model.WhoAmI{}, err
	}

	return whoami, nil
}

func (m Monzo) Accounts() (model.Accounts, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)

	request := httpc.HTTP{
		TargetURL: GetURL(AccountsURL),
		Method:    http.MethodGet,
		Headers:   headers,
		Form:      nil,
	}

	var accounts model.Accounts
	if err := request.JSON(&accounts); err != nil {
		return model.Accounts{}, err
	}

	return accounts, nil
}

func (m Monzo) CurrentAccount() (model.Account, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)

	request := httpc.HTTP{
		TargetURL: GetURL(AccountsURL),
		Method:    http.MethodGet,
		Headers:   headers,
		Form:      nil,
	}

	var accounts model.Accounts
	if err := request.JSON(&accounts); err != nil {
		return model.Account{}, err
	}

	return accounts.Array[0], nil
}

func (m Monzo) Balance(accountID string) (model.Balance, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)

	targetURL := GetURL(BalanceURL, accountID)
	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodGet,
		Headers:   headers,
		Form:      nil,
	}

	var balance model.Balance
	if err := request.JSON(&balance); err != nil {
		return model.Balance{}, err
	}

	return balance, nil
}

func (m Monzo) Pots() (model.Pots, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)

	request := httpc.HTTP{
		TargetURL: GetURL(PotsURL),
		Method:    http.MethodGet,
		Headers:   headers,
		Form:      nil,
	}

	var pots model.Pots

	if err := request.JSON(&pots); err != nil {
		return model.Pots{}, err
	}

	return pots, nil
}

func (m Monzo) RegisterWebhook(accountID string) (model.Webhook, error) {
	headers := httpc.Headers{}
	headers.FormURLEncoded()
	headers.Authorization(m.tokenType, m.accessToken)

	data := map[string]string{
		"account_id": accountID,
		"url":        WebhookURI,
	}

	request := httpc.HTTP{
		TargetURL: GetURL(WebhookCreateURL),
		Method:    http.MethodPost,
		Headers:   headers,
		Form:      data,
	}

	webhookBody := model.WebhookBody{}
	if err := request.JSON(&webhookBody); err != nil {
		return model.Webhook{}, err
	}

	return webhookBody.Webhook, nil
}

func (m Monzo) DeleteWebhook(webhookID string) error {
	headers := httpc.Headers{}
	headers.Authorization(m.tokenType, m.accessToken)

	request := httpc.HTTP{
		TargetURL: GetURL(WebhookDeleteURL, webhookID),
		Method:    http.MethodDelete,
		Headers:   headers,
		Form:      nil,
	}

	if _, err := request.String(); err != nil {
		return err
	}

	return nil
}

func (m Monzo) Webhooks(accountID string) ([]model.Webhook, error) {
	targetURL := GetURL(WebhookGetURL, accountID)

	headers := httpc.Headers{}
	headers.Authorization(m.tokenType, m.accessToken)

	request := httpc.HTTP{
		TargetURL: targetURL,
		Headers:   headers,
		Form:      nil,
	}

	var webhooks model.Webhooks
	if err := request.JSON(&webhooks); err != nil {
		logio.Println(err)
		return nil, err
	}

	return webhooks.Array, nil
}

func (m Monzo) Withdraw(sourcePotID string, destinationAccountID string, amount int64) (model.Pot, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)
	headers.FormURLEncoded()

	data := make(map[string]string)
	data["destination_account_id"] = destinationAccountID
	data["amount"] = strconv.FormatInt(amount, 10)
	data["dedupe_id"] = uuid.Token()

	targetURL := GetURL(WithdrawURL, sourcePotID)

	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodPut,
		Headers:   headers,
		Form:      data,
	}

	var pot model.Pot
	if err := request.JSON(&pot); err != nil {
		return model.Pot{}, err
	}

	return pot, nil
}

func (m Monzo) Deposit(targetPotID string, sourceAccountID string, amount int64) (model.Pot, error) {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)
	headers.FormURLEncoded()

	data := make(map[string]string)
	data["source_account_id"] = sourceAccountID
	data["amount"] = strconv.FormatInt(amount, 10)
	data["dedupe_id"] = uuid.Token()

	targetURL := GetURL(DepositURL, targetPotID)

	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodPut,
		Headers:   headers,
		Form:      data,
	}

	var pot model.Pot
	if err := request.JSON(&pot); err != nil {
		return model.Pot{}, err
	}

	return pot, nil
}

func (m Monzo) CreateFeedItem(accountID string, title string, body string, imageURL string) error {
	headers := make(httpc.Headers)
	headers.Authorization(m.tokenType, m.accessToken)
	headers.FormURLEncoded()

	data := make(map[string]string)
	data["account_id"] = accountID
	data["type"] = "basic"
	// data["url"] = "https://monzo.com"
	data["params[title]"] = title
	data["params[body]"] = body
	data["params[image_url]"] = imageURL

	targetURL := GetURL(FeedItemCreateURL)

	request := httpc.HTTP{
		TargetURL: targetURL,
		Method:    http.MethodPost,
		Headers:   headers,
		Form:      data,
	}

	if _, err := request.Status(); err != nil {
		return err
	}

	return nil
}
