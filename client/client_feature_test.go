// +build feature

package client_test

import (
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

func newClient() *client.Client {
	config := client.Config{
		URL:         os.Getenv("GOHBCI_URL"),
		AccountID:   os.Getenv("GOHBCI_USERID"),
		BankID:      os.Getenv("GOHBCI_BLZ"),
		PIN:         os.Getenv("GOHBCI_PIN"),
		HBCIVersion: client.Version220,
	}
	testAccount = domain.AccountConnection{AccountID: config.AccountID, CountryCode: 280, BankID: config.BankID}
	c, err := client.New(config)
	if err != nil {
		panic(err)
	}
	return c
}
