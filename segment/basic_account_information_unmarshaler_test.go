package segment

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestAccountInformationSegmentUnmarshalHBCI(t *testing.T) {
	testCases := []struct {
		desc                   string
		rawSegment             string
		expectedSegmentBuilder func() *AccountInformationSegment
		expectedError          error
	}{
		{
			desc:       "fully specced version 4 segment",
			rawSegment: "HIUPD:1:4:4+123456::280:10000000+12345+EUR+Muster+Max+Sichteinlagen++DKPAE:1'",
			expectedSegmentBuilder: func() *AccountInformationSegment {
				v4 := &AccountInformationV4{
					AccountConnection: element.NewAccountConnection(domain.AccountConnection{AccountID: "123456", CountryCode: 280, BankID: "10000000"}),
					UserID:            element.NewIdentification("12345"),
					AccountCurrency:   element.NewCurrency("EUR"),
					Name1:             element.NewAlphaNumeric("Muster", 27),
					Name2:             element.NewAlphaNumeric("Max", 27),
					AccountProductID:  element.NewAlphaNumeric("Sichteinlagen", 30),
					AllowedBusinessTransactions: element.NewAllowedBusinessTransactions(
						domain.BusinessTransaction{ID: "DKPAE", NeededSignatures: 1},
					),
				}
				v4.Segment = NewReferencingBasicSegment(1, 4, v4)
				return &AccountInformationSegment{v4}
			},
		},
		{
			desc:       "barebone version 5 segment",
			rawSegment: "HIUPD:1:5:4++Login+++Name++++HKPSA:1'",
			expectedSegmentBuilder: func() *AccountInformationSegment {
				v5 := &AccountInformationV5{
					UserID: element.NewIdentification("Login"),
					Name1:  element.NewAlphaNumeric("Name", 27),
					AllowedBusinessTransactions: element.NewAllowedBusinessTransactions(
						domain.BusinessTransaction{ID: "HKPSA", NeededSignatures: 1},
					),
				}
				v5.Segment = NewReferencingBasicSegment(1, 4, v5)
				return &AccountInformationSegment{v5}
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			account := &AccountInformationSegment{}

			err := account.UnmarshalHBCI([]byte(tt.rawSegment))

			if err != nil {
				t.Logf("Expected no error, got %T:%v\n", err, err)
				t.Fail()
			}

			expected := tt.expectedSegmentBuilder()

			compareAccountInformationSegments(t, expected, account)
		})
	}
}

func compareAccountInformationSegments(t *testing.T, expected, actual *AccountInformationSegment) {
	t.Helper()
	var expectedString, actualString string = "<nil>", "<nil>"
	if expected != nil {
		expectedString = expected.String()
	}
	if actual != nil {
		actualString = actual.String()
	}
	if expectedString != actualString {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expectedString, actualString)
		t.Fail()
	}
}
