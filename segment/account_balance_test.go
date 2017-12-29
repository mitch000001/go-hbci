package segment

import (
	"testing"
	"time"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestAccountBalanceResponseSegmentUnmarshalHBCI(t *testing.T) {
	test := "HISAL:4:5:3+100000000::280:10000000+Sichteinlagen+EUR+C:1000,15:EUR:20150812+C:20,:EUR:20150812+500,:EUR+1499,85:EUR'"
	date := time.Date(2015, 8, 12, 0, 0, 0, 0, time.Local)

	expectedSegment := &AccountBalanceResponseSegment{
		AccountConnection:  element.NewAccountConnection(domain.AccountConnection{AccountID: "100000000", CountryCode: 280, BankID: "10000000"}),
		AccountProductName: element.NewAlphaNumeric("Sichteinlagen", 35),
		AccountCurrency:    element.NewCurrency("EUR"),
		BookedBalance:      element.NewBalance(domain.Amount{Amount: 1000.15, Currency: "EUR"}, date, false),
		EarmarkedBalance:   element.NewBalance(domain.Amount{Amount: 20, Currency: "EUR"}, date, false),
		CreditLimit:        element.NewAmount(500, "EUR"),
		AvailableAmount:    element.NewAmount(1499.85, "EUR"),
	}
	expectedSegment.Segment = NewReferencingBasicSegment(4, 3, expectedSegment)

	expected := expectedSegment.String()

	segment := &AccountBalanceResponseSegment{}

	err := segment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected error to be nil, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := segment.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
