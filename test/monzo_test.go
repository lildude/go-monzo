package test

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"time"

	"encoding/json"

	"strings"

	"github.com/gurparit/go-monzo/monzo"
	"github.com/gurparit/go-monzo/monzo/model"
)

func TestMonzoAuthHeaderInvalid(t *testing.T) {
	authHeaderInvalid := `
{
	"forbidden": "authorization header invalid"
}
`

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", "Bearer x-access-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(authHeaderInvalid))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.AccountsURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Accounts()
	if err != nil && !strings.Contains(err.Error(), authHeaderInvalid) {
		t.Log(err)
		t.FailNow()
	}
}

func TestMonzoAuthHeaderMissing(t *testing.T) {
	authHeaderMissing := `
{
	"error": "authorization header missing"
}
`

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", "Bearer x-access-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(authHeaderMissing))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.AccountsURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Accounts()
	if err != nil && !strings.Contains(err.Error(), authHeaderMissing) {
		t.Log(err)
		t.FailNow()
	}
}

func TestMonzoAuthenticationExpired(t *testing.T) {
	badAccessToken := `
{
	"error": "{\"code\":\"unauthorized.bad_access_token.expired\",\"error\":\"invalid_token\",\"error_description\":\"Access token has expired\",\"message\":\"Access token has expired\"}\n"
}
`

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", "Bearer x-access-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(badAccessToken))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.AccountsURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Accounts()
	if err != nil && !strings.Contains(err.Error(), badAccessToken) {
		t.Log(err)
		t.FailNow()
	}
}

func TestMonzoCallback(t *testing.T) {
	timeNow := time.Now().Add(time.Second * 21600)
	expectedNewUser := model.User{
		ID:           2,
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "new-x-access-token",
		RefreshToken: "new-x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Content-Type", "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		IsEqual(t, "Method", http.MethodPost, r.Method)

		// Test Form
		IsEqual(t, "grant_type", "authorization_code", r.PostFormValue("grant_type"))
		IsEqual(t, "code", "x-code", r.PostFormValue("code"))

		data := map[string]interface{}{
			"user_id":       expectedNewUser.UserID,
			"client_id":     expectedNewUser.ClientID,
			"access_token":  expectedNewUser.AccessToken,
			"refresh_token": expectedNewUser.RefreshToken,
			"token_type":    expectedNewUser.TokenType,
			"expires_in":    expectedNewUser.ExpiresIn,
		}

		response, err := json.Marshal(data)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.Oauth2URL, testHttp.URL)

	testUser, err := monzo.Callback("x-code")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "user.id", expectedNewUser.UserID, testUser.UserID)
	IsEqual(t, "user.client_id", expectedNewUser.ClientID, testUser.ClientID)
	IsEqual(t, "user.access_token", expectedNewUser.AccessToken, testUser.AccessToken)
	IsEqual(t, "user.refresh_token", expectedNewUser.RefreshToken, testUser.RefreshToken)
	IsEqual(t, "user.token_type", expectedNewUser.TokenType, testUser.TokenType)
	IsEqual(t, "user.expires_in", expectedNewUser.ExpiresIn, testUser.ExpiresIn)
}

func TestMonzoRefresh(t *testing.T) {
	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		ID:           1,
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	// Expected New User
	expectedNewUser := model.User{
		ID:           expectedUser.ID,
		UserID:       expectedUser.UserID,
		ClientID:     expectedUser.ClientID,
		AccessToken:  "new-x-access-token",
		RefreshToken: "new-x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Content-Type", r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		IsEqual(t, "Method", http.MethodPost, r.Method)

		// Test Form
		IsEqual(t, "grant_type", "refresh_token", r.PostFormValue("grant_type"))
		IsEqual(t, "refresh_token", "x-refresh-token", r.PostFormValue("refresh_token"))

		data := map[string]interface{}{
			"user_id":       expectedNewUser.UserID,
			"client_id":     expectedNewUser.ClientID,
			"access_token":  expectedNewUser.AccessToken,
			"refresh_token": expectedNewUser.RefreshToken,
			"token_type":    expectedNewUser.TokenType,
			"expires_in":    expectedNewUser.ExpiresIn,
		}

		response, err := json.Marshal(data)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.Oauth2URL, testHttp.URL)

	testUser, err := monzo.Refresh("x-refresh-token")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "user.id", expectedNewUser.UserID, testUser.UserID)
	IsEqual(t, "user.client_id", expectedNewUser.ClientID, testUser.ClientID)
	IsEqual(t, "user.access_token", expectedNewUser.AccessToken, testUser.AccessToken)
	IsEqual(t, "user.refresh_token", expectedNewUser.RefreshToken, testUser.RefreshToken)
	IsEqual(t, "user.token_type", expectedNewUser.TokenType, testUser.TokenType)
	IsEqual(t, "user.expires_in", expectedNewUser.ExpiresIn, testUser.ExpiresIn)
}

func TestMonzoAccountsNew(t *testing.T) {
	// Expected Account
	expectedAccount := model.Account{
		AccountID:   "x-account-id",
		Description: "My Current Account",
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", "Bearer x-access-token", r.Header.Get("Authorization"))

		data := map[string]interface{}{
			"accounts": []map[string]interface{}{{
				"id":          expectedAccount.AccountID,
				"description": expectedAccount.Description,
				"created":     "2018-01-01T12:12:21Z",
			}},
		}

		response, err := json.Marshal(data)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.AccountsURL, testHttp.URL)

	testMonzoAccounts, err := monzo.New("Bearer", "x-access-token").Accounts()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "count(accounts)", 1, len(testMonzoAccounts.Array))

	testAccount := testMonzoAccounts.Array[0]
	IsEqual(t, "account", expectedAccount, testAccount)
}

func TestMonzoAccountsUpdate(t *testing.T) {
	// Expected Account
	expectedAccount := model.Account{
		AccountID:   "x-account-id",
		Description: "My Current Account",
	}

	// Expected New Account
	expectedNewAccount := model.Account{
		ID:          expectedAccount.ID,
		UserID:      expectedAccount.UserID,
		AccountID:   "new-x-account-id",
		Description: "My New Current Account",
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", "Bearer x-access-token", r.Header.Get("Authorization"))

		data := map[string]interface{}{
			"accounts": []map[string]interface{}{{
				"id":          expectedNewAccount.AccountID,
				"description": expectedNewAccount.Description,
				"created":     "2018-01-01T12:12:21Z",
			}},
		}

		response, err := json.Marshal(data)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.AccountsURL, testHttp.URL)

	testMonzoAccounts, err := monzo.New("Bearer", "x-access-token").Accounts()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "count(accounts)", 1, len(testMonzoAccounts.Array))

	testAccount := testMonzoAccounts.Array[0]
	IsEqual(t, "account", expectedNewAccount, testAccount)
}

func IsEqual(t *testing.T, key string, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Logf("for %s;", key)
		t.Logf("expected %s;", expected)
		t.Logf("actual: %s;", actual)
		t.FailNow()
	}
}
