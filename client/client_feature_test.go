// +build feature

package client_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
)

var testAccount domain.AccountConnection

func TestClientAccountTransactions(t *testing.T) {
	c := newClient()

	timeframe := domain.Timeframe{
		StartDate: domain.NewShortDate(time.Now().AddDate(0, 0, -10)),
	}
	transactions, err := c.AccountTransactions(testAccount, timeframe, false, "")

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}

	if transactions == nil {
		t.Logf("Expected transactions not to be nil\n")
		t.Fail()
	}

	for _, tr := range transactions {
		fmt.Printf("Transaction: %s\n", tr)
	}
}

func TestPinTanDialogAccountInformation(t *testing.T) {
	t.Skip()
	c := newClient()

	err := c.AccountInformation(testAccount, true)

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}
}

func TestPinTanDialogBalances(t *testing.T) {
	t.Skip()
	c := newClient()

	balances, err := c.AccountBalances(testAccount, true)

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}

	if balances == nil {
		t.Logf("Expected balances not to be nil\n")
		t.Fail()
	}
}

func TestClientAccounts(t *testing.T) {
	c := newClient()

	accounts, err := c.Accounts()

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}

	if accounts == nil {
		t.Logf("Expected accounts not to be nil\n")
		t.Fail()
	}

	for _, account := range accounts {
		fmt.Printf("Account: %s\nProduct name: %s\nCurrency: %s\n", account.AccountConnection.AccountID, account.ProductID, account.Currency)
	}
}

func newClient() *client.Client {
	configFile, err := os.Open("../.fints300_haspa.json")
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	var config client.Config
	if configFile != nil {
		jsonDecoder := json.NewDecoder(configFile)
		err = jsonDecoder.Decode(&config)
		if err != nil {
			panic(err)
		}
	} else {
		config = client.Config{
			URL:         os.Getenv("GOHBCI_URL"),
			AccountID:   os.Getenv("GOHBCI_USERID"),
			BankID:      os.Getenv("GOHBCI_BLZ"),
			PIN:         os.Getenv("GOHBCI_PIN"),
			HBCIVersion: domain.HBCIVersion220,
		}
	}
	testAccount = domain.AccountConnection{AccountID: config.AccountID, CountryCode: 280, BankID: config.BankID}
	c, err := client.New(config)
	if err != nil {
		panic(err)
	}
	return c
}
