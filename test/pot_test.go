package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"strings"

	"github.com/gurparit/go-monzo/model"
	"github.com/gurparit/go-monzo/monzo"
)

func TestPotWithdrawSuccess(t *testing.T) {
	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
		Created:     time.Now(),
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Destination Account ID", expectedAccount.ID, r.PostFormValue("destination_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		response, err := json.Marshal(expectedPot)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.WithdrawURL, testHttp.URL)

	pot, err := monzo.New("Bearer", "x-access-token").Withdraw(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "pot.id", expectedPot.ID, pot.ID)
	IsEqual(t, "pot.name", expectedPot.Name, pot.Name)
	IsEqual(t, "pot.balance", expectedPot.Balance, pot.Balance)
	IsEqual(t, "pot.currency", expectedPot.Currency, pot.Currency)
	IsEqual(t, "pot.deleted", expectedPot.Deleted, pot.Deleted)
}

func TestPotDepositSuccess(t *testing.T) {
	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
		Created:     time.Now(),
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Source Account ID", expectedAccount.ID, r.PostFormValue("source_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		response, err := json.Marshal(expectedPot)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.DepositURL, testHttp.URL)

	pot, err := monzo.New("Bearer", "x-access-token").Deposit(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "pot.id", expectedPot.ID, pot.ID)
	IsEqual(t, "pot.name", expectedPot.Name, pot.Name)
	IsEqual(t, "pot.balance", expectedPot.Balance, pot.Balance)
	IsEqual(t, "pot.currency", expectedPot.Currency, pot.Currency)
	IsEqual(t, "pot.deleted", expectedPot.Deleted, pot.Deleted)
}

func TestPotWithdrawDeletedFail(t *testing.T) {
	sampleDepositFail := `
{
	"error": "cannot access deleted pots"
}
`

	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
		Created:     time.Now(),
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Destination Account ID", expectedAccount.ID, r.PostFormValue("destination_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(sampleDepositFail))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.WithdrawURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Withdraw(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil && !strings.Contains(err.Error(), "cannot access deleted pots") {
		t.Fail()
	}
}

func TestPotDepositDeletedFail(t *testing.T) {
	sampleDepositFail := `
{
	"error": "cannot access deleted pots"
}
`

	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
		Created:     time.Now(),
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Source Account ID", expectedAccount.ID, r.PostFormValue("source_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(sampleDepositFail))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.DepositURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Deposit(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil && !strings.Contains(err.Error(), "cannot access deleted pots") {
		t.Fail()
	}
}

func TestPotWithdrawInsufficientFundsFail(t *testing.T) {
	sampleDepositFail := `
{
	"error": "cannot withdraw amount, not enough money in pot"
}
`

	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Destination Account ID", expectedAccount.ID, r.PostFormValue("destination_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(sampleDepositFail))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.WithdrawURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Withdraw(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil && !strings.Contains(err.Error(), "cannot withdraw amount, not enough money in pot") {
		t.Fail()
	}
}

func TestPotDepositInsufficientFundsFail(t *testing.T) {
	sampleDepositFail := `
{
	"error": "{\"code\":\"bad_request.insufficient_funds\",\"message\":\"You can't deposit more than your current account balance\"}\n"
}
`

	// Expected User
	timeNow := time.Now().Add(time.Second * 21600)
	expectedUser := model.User{
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedAccount := model.Account{
		ID:          "x-account-id",
		Description: "My Current Account",
		Created:     time.Now(),
	}

	expectedPot := model.Pot{
		ID:       "pot_00009exampleP0tOxWb",
		Name:     "Flying Lessons",
		Balance:  350000,
		Currency: "GBP",
		Deleted:  false,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodPut, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))
		IsEqual(t, "Source Account ID", expectedAccount.ID, r.PostFormValue("source_account_id"))
		IsEqual(t, "Amount", "5000", r.PostFormValue("amount"))
		IsEqual(t, "Dedupe ID", true, r.PostFormValue("dedupe_id") != "")

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(sampleDepositFail))
	}))

	defer testHttp.Close()

	monzo.SetURL(monzo.DepositURL, testHttp.URL)

	_, err := monzo.New("Bearer", "x-access-token").Deposit(expectedPot.ID, expectedAccount.ID, 5000)
	if err != nil && !strings.Contains(err.Error(), "{\"code\":\"bad_request.insufficient_funds\",\"message\":\"You can't deposit more than your current account balance\"}\n") {
		t.Fail()
	}
}
