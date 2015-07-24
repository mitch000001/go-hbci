package segment

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestCommonUserParameterDataSegmentUnmarshalHBCI(t *testing.T) {
	test := "HIUPA:5:2:7+12345+4+0'"

	segment := &CommonUserParameterDataSegment{}

	err := segment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := &CommonUserParameterDataSegment{
		UserID:     element.NewIdentification("12345"),
		UPDVersion: element.NewNumber(4, 3),
		UPDUsage:   element.NewNumber(0, 1),
	}
	expected.Segment = NewReferencingBasicSegment(5, 7, expected)

	expectedString := expected.String()
	actualString := segment.String()

	if expectedString != actualString {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expectedString, actualString)
		t.Fail()
	}
}

func TestAccountInformationSegmentUnmarshalHBCI(t *testing.T) {
	test := "HIUPD:1:4:4+123456::280:10000000+12345+EUR+Muster+Max+Sichteinlagen++DKPAE:1"

	account := &AccountInformationSegment{}

	err := account.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := &AccountInformationSegment{
		AccountConnection:           element.NewAccountConnection(domain.AccountConnection{AccountID: "123456", CountryCode: 280, BankID: "10000000"}),
		UserID:                      element.NewIdentification("12345"),
		AccountCurrency:             element.NewCurrency("EUR"),
		Name1:                       element.NewAlphaNumeric("Muster", 27),
		Name2:                       element.NewAlphaNumeric("Max", 27),
		AccountProductID:            element.NewAlphaNumeric("Sichteinlagen", 30),
		AllowedBusinessTransactions: element.NewAllowedBusinessTransactions(domain.BusinessTransaction{ID: "DKPAE", NeededSignatures: 1}),
	}
	expected.Segment = NewReferencingBasicSegment(1, 4, expected)

	expectedString := expected.String()
	actualString := account.String()

	if expectedString != actualString {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expectedString, actualString)
		t.Fail()
	}
}
