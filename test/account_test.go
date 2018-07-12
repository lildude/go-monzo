package test

import (
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"
	"time"

	"fmt"

	"github.com/gurparit/go-monzo/model"
	"github.com/gurparit/go-monzo/monzo"
)

func TestAccountBalanceNew(t *testing.T) {
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

	testBalance, err := monzo.New("Bearer", "x-access-token").Balance("x-account-id")
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

	testBalance, err := monzo.New("Bearer", "x-access-token").Balance("x-account-id")
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
		UserID:       "x-user-id",
		ClientID:     "x-client-id",
		AccessToken:  "x-access-token",
		RefreshToken: "x-refresh-token",
		ExpiryDate:   timeNow,
		TokenType:    "Bearer",
		ExpiresIn:    21600,
	}

	expectedActivePot := model.Pot{
		ID:       "x-active-pot-id",
		Name:     "x-active-pot-name",
		Balance:  4000,
		Currency: "GBP",
		Deleted:  false,
	}

	expectedDeletedPot := model.Pot{
		ID:       "x-deleted-pot-id",
		Name:     "x-deleted-pot-name",
		Balance:  4000,
		Currency: "GBP",
		Deleted:  true,
	}

	testHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test Headers
		authHeader := fmt.Sprintf("%s %s", expectedUser.TokenType, expectedUser.AccessToken)

		IsEqual(t, "Method", http.MethodGet, r.Method)
		IsEqual(t, "Authorization", authHeader, r.Header.Get("Authorization"))

		data := model.Monzo{
			Pots: []model.Pot{
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

	testMonzo, err := monzo.New("Bearer", "x-access-token").Pots()

	IsEqual(t, "error", nil, err)

	actualActivePot := testMonzo.Pots[0]
	actualDeletedPot := testMonzo.Pots[1]

	IsEqual(t, "len(pots)", 2, len(testMonzo.Pots))

	// Active Pot
	IsEqual(t, "active pot.id", expectedActivePot.ID, actualActivePot.ID)
	IsEqual(t, "active pot.name", expectedActivePot.Name, actualActivePot.Name)
	IsEqual(t, "active pot.balance", expectedActivePot.Balance, actualActivePot.Balance)
	IsEqual(t, "active pot.currency", expectedActivePot.Currency, actualActivePot.Currency)
	IsEqual(t, "active pot.deleted", expectedActivePot.Deleted, actualActivePot.Deleted)

	// Deleted Pot
	IsEqual(t, "deleted pot.id", expectedDeletedPot.ID, actualDeletedPot.ID)
	IsEqual(t, "deleted pot.name", expectedDeletedPot.Name, actualDeletedPot.Name)
	IsEqual(t, "deleted pot.balance", expectedDeletedPot.Balance, actualDeletedPot.Balance)
	IsEqual(t, "deleted pot.currency", expectedDeletedPot.Currency, actualDeletedPot.Currency)
	IsEqual(t, "deleted pot.deleted", expectedDeletedPot.Deleted, actualDeletedPot.Deleted)
}
