package client

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	https "github.com/mitch000001/go-hbci/transport/https"
)

func TestClientBalances(t *testing.T) {
	transport := &https.MockHTTPTransport{}
	defer setMockHTTPTransport(transport)()

	c := newTestClient()

	account := domain.AccountInformation{
		AccountConnection: domain.AccountConnection{AccountID: "100000000", CountryCode: 280, BankID: "10000000"},
		UserID:            "100000000",
		Currency:          "EUR",
		Name1:             "Muster",
		Name2:             "Max",
		AllowedBusinessTransactions: []domain.BusinessTransaction{
			{ID: "HKSAL", NeededSignatures: 1},
		},
	}

	c.pinTanDialog.Accounts = []domain.AccountInformation{
		account,
	}

	syncResponse := encryptedTestMessage(
		"abcde",
		"HIRMG:2:2:1+0020::Auftrag entgegengenommen'",
		"HISYN:193:4:5+LRZYhZNbV2IBAAAd0?+VNqlkXrAQA'",
		"HIBPA:2:2:+12+280:10000000+Bank Name+3+1+201:210:220+0'",
		"HISALS:3:5:4+3+1'",
	)
	initResponse := encryptedTestMessage(
		"abcde",
		"HIRMG:2:2:1+0020::Auftrag entgegengenommen'",
		"HIKIM:3:2+ec-Karte+Ihre neue ec-Karte liegt zur Abholung bereit.'",
	)
	balanceResponse := encryptedTestMessage(
		"abcde",
		"HIRMG:2:2:1+0020::Auftrag entgegengenommen'",
		"HISAL:3:5:1+100000000::280:10000000+Sichteinlagen+EUR+C:1000,15:EUR:20150812+C:20,:EUR:20150812+500,:EUR+1499,85:EUR'",
	)
	dialogEndResponseMessage := encryptedTestMessage("abcde", "HIRMG:2:2:1+0020::Der Auftrag wurde ausgeführt'")

	transport.SetResponsePayloads([][]byte{
		syncResponse,
		dialogEndResponseMessage,
		initResponse,
		balanceResponse,
		dialogEndResponseMessage,
	})

	balances, err := c.AccountBalances(account.AccountConnection, true)
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	date, _ := time.Parse("20060102", "20150812")

	expectedBalance := domain.AccountBalance{
		Account:     domain.AccountConnection{AccountID: "100000000", CountryCode: 280, BankID: "10000000"},
		ProductName: "Sichteinlagen",
		Currency:    "EUR",
		BookedBalance: domain.Balance{
			Amount:           domain.Amount{Amount: 1000.15, Currency: "EUR"},
			TransmissionDate: date,
		},
		EarmarkedBalance: &domain.Balance{
			Amount:           domain.Amount{Amount: 20, Currency: "EUR"},
			TransmissionDate: date,
		},
		CreditLimit:     &domain.Amount{Amount: 500, Currency: "EUR"},
		AvailableAmount: &domain.Amount{Amount: 1499.85, Currency: "EUR"},
	}

	if len(balances) != 1 {
		t.Logf("Expected balances length to equal 1, was %d\n", len(balances))
		t.Fail()
	} else {
		if !reflect.DeepEqual(balances[0], expectedBalance) {
			t.Logf("Expected balance to equal\n%#v\n\tgot\n%#v\n", expectedBalance, balances[0])
			t.Fail()
		}
	}
}

func newTestClient() *Client {
	config := Config{
		URL:         "https://localhost",
		AccountID:   "12345",
		BankID:      "10000000",
		PIN:         "abcde",
		HBCIVersion: domain.HBCIVersion220,
	}
	c, err := New(config)
	if err != nil {
		panic(err)
	}
	c.pinTanDialog.SetClientSystemID("xyz")
	return c
}

func setMockHTTPTransport(transport http.RoundTripper) func() {
	originHTTPTransport := http.DefaultTransport
	http.DefaultTransport = transport
	return func() {
		http.DefaultTransport = originHTTPTransport
	}
}
