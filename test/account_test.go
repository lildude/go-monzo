package test

import (
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	"fmt"

	"github.com/gurparit/go-monzo/monzo"
	"github.com/gurparit/go-monzo/monzo/model"
)

func TestAccountBalanceNew(t *testing.T) {
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

	expectedNewBalance := model.Balance{
		Balance:      12000,
		TotalBalance: 22800,
		Currency:     "GBP",
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))

		data := map[string]interface{}{
			"spend_today":   0,
			"balance":       expectedNewBalance.Balance,
			"total_balance": expectedNewBalance.TotalBalance,
			"currency":      expectedNewBalance.Currency,
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

	monzo.SetURL(monzo.BalanceURL, testHttp.URL)

	testBalance, err := monzo.Balance("Bearer", "x-access-token", "x-account-id")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "balance", expectedNewBalance, testBalance)
}

func TestAccountBalanceUpdate(t *testing.T) {
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

	expectedBalance := model.Balance{
		Balance:      12000,
		TotalBalance: 22800,
		Currency:     "GBP",
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))

		data := map[string]interface{}{
			"spend_today":   0,
			"balance":       expectedBalance.Balance,
			"total_balance": expectedBalance.TotalBalance,
			"currency":      expectedBalance.Currency,
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

	monzo.SetURL(monzo.BalanceURL, testHttp.URL)

	testBalance, err := monzo.Balance("Bearer", "x-access-token", "x-account-id")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	IsEqual(t, "balance", expectedBalance, testBalance)
}

func TestAccountPots(t *testing.T) {
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

	expectedAccount := model.Account{
		ID:          1,
		UserID:      expectedUser.UserID,
		AccountID:   "x-account-id",
		Description: "My Current Account",
	}

	expectedActivePot := model.Pot{
		ID:        0,
		PotID:     "x-active-pot-id",
		Name:      "x-active-pot-name",
		AccountID: expectedAccount.AccountID,
		Balance:   4000,
		Currency:  "GBP",
		Deleted:   false,
	}

	expectedDeletedPot := model.Pot{
		ID:        0,
		PotID:     "x-deleted-pot-id",
		Name:      "x-deleted-pot-name",
		AccountID: expectedAccount.AccountID,
		Balance:   4000,
		Currency:  "GBP",
		Deleted:   true,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))

		data := model.Pots{
			Array: []model.Pot{
				expectedActivePot,
				expectedDeletedPot,
			},
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

	monzo.SetURL(monzo.PotsURL, testHttp.URL)

	pots, err := monzo.Pots("Bearer", "x-access-token")

	IsEqual(t, "error", nil, err)

	actualActivePot := pots.Array[0]
	actualDeletedPot := pots.Array[1]

	IsEqual(t, "len(pots)", 2, len(pots.Array))

	// Active Pot
	IsEqual(t, "active pot.id", expectedActivePot.ID, actualActivePot.ID)
	IsEqual(t, "active pot.name", expectedActivePot.Name, actualActivePot.Name)
	IsEqual(t, "active pot.balance", expectedActivePot.Balance, actualActivePot.Balance)
	IsEqual(t, "active pot.currency", expectedActivePot.Currency, actualActivePot.Currency)
	IsEqual(t, "active pot.account_id", expectedActivePot.AccountID, actualActivePot.AccountID)
	IsEqual(t, "active pot.deleted", expectedActivePot.Deleted, actualActivePot.Deleted)

	// Deleted Pot
	IsEqual(t, "deleted pot.id", expectedDeletedPot.ID, actualDeletedPot.ID)
	IsEqual(t, "deleted pot.name", expectedDeletedPot.Name, actualDeletedPot.Name)
	IsEqual(t, "deleted pot.balance", expectedDeletedPot.Balance, actualDeletedPot.Balance)
	IsEqual(t, "deleted pot.currency", expectedDeletedPot.Currency, actualDeletedPot.Currency)
	IsEqual(t, "deleted pot.account_id", expectedDeletedPot.AccountID, actualDeletedPot.AccountID)
	IsEqual(t, "deleted pot.deleted", expectedDeletedPot.Deleted, actualDeletedPot.Deleted)
}
