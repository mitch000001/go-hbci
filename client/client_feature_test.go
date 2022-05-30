package client_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mitch000001/go-hbci/client"
	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/iban"
)

var testAccount domain.AccountConnection
var sepaTestAccount domain.InternationalAccountConnection

func TestClientAccountTransactions(t *testing.T) {
	skipWhenE2EDisabled(t)
	c := newClient()

	timeframe := domain.Timeframe{
		StartDate: domain.NewShortDate(time.Now().AddDate(0, 0, -90)),
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

func TestClientSepaAccountTransactions(t *testing.T) {
	skipWhenE2EDisabled(t)
	c := newClient()

	timeframe := domain.Timeframe{
		StartDate: domain.NewShortDate(time.Now().AddDate(0, 0, -90)),
	}
	transactions, err := c.SepaAccountTransactions(sepaTestAccount, timeframe, false, "")

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

func TestClientAccountInformation(t *testing.T) {
	skipWhenE2EDisabled(t)
	c := newClient()

	err := c.AccountInformation(testAccount, true)

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		// t.Fail()
	}
}

func TestClientAccountBalances(t *testing.T) {
	skipWhenE2EDisabled(t)
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

	for _, balance := range balances {
		t.Logf("Balance: %s\n", balance)
	}
}

func TestClientAccounts(t *testing.T) {
	skipWhenE2EDisabled(t)
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
		t.Logf("Account: %s\tProduct name: %s\tCurrency: %s\n", account.AccountConnection.AccountID, account.ProductID, account.Currency)
		t.Logf("Allowed transactions:\n")
		for _, allowedTransaction := range account.AllowedBusinessTransactions {
			t.Logf("ID: %s\tNeeded signatures: %d\n", allowedTransaction.ID, allowedTransaction.NeededSignatures)
		}
	}
}

func TestClientStatus(t *testing.T) {
	skipWhenE2EDisabled(t)
	c := newClient()

	statuus, err := c.Status(time.Now().Add(-48*time.Hour), time.Now(), 10, "")

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}

	for _, status := range statuus {
		t.Logf("Status: %s\n", status)
	}
}

func TestAnonymousClientCommunicationAccess(t *testing.T) {
	skipWhenE2EDisabled(t)
	a := &client.AnonymousClient{
		Client: newClient(),
	}

	from := domain.BankID{280, "78050000"}
	to := domain.BankID{280, "78050000"}

	res, err := a.CommunicationAccess(from, to, 10)

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		// t.Fail()
	}

	if res != nil {
		t.Logf("Response: %s\n", res)
	}
}

func newClient() *client.Client {
	configFile, err := os.Open("../.fints300.json")
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
			URL:       os.Getenv("GOHBCI_URL"),
			AccountID: os.Getenv("GOHBCI_USERID"),
			BankID:    os.Getenv("GOHBCI_BLZ"),
			PIN:       os.Getenv("GOHBCI_PIN"),
		}
	}
	config.StoreBankParameterData = os.Getenv("GOHBCI_STORE_BPD") == "true"
	config.StoreUserParameterData = os.Getenv("GOHBCI_STORE_UPD") == "true"
	testAccount = domain.AccountConnection{AccountID: config.AccountID, CountryCode: 280, BankID: config.BankID}
	i, err := iban.NewGerman(config.BankID, config.AccountID)
	if err != nil {
		panic(err)
	}
	sepaTestAccount = domain.InternationalAccountConnection{
		IBAN:      string(i),
		AccountID: config.AccountID,
		BankID:    domain.BankID{CountryCode: 280, ID: config.BankID},
	}
	c, err := client.New(config)
	if err != nil {
		panic(err)
	}
	return c
}

func skipWhenE2EDisabled(t *testing.T) {
	if testE2E() {
		return
	}
	t.Skip("Enable client e2e tests by setting env var `GOHBCI_TEST_E2E_CLIENT` to `true`")
}

func testE2E() bool {
	val, ok := os.LookupEnv("GOHBCI_TEST_E2E_CLIENT")
	if !ok {
		return false
	}
	testE2E, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return testE2E
}
