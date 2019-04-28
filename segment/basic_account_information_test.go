package segment

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestAccountInformationSegment_Account(t *testing.T) {
	testCases := []struct {
		desc            string
		segment         AccountInformation
		expectedAccount domain.AccountInformation
	}{
		{
			desc: "v4 all mandatory fields set",
			segment: &AccountInformationV4{
				AccountConnection: element.NewAccountConnection(
					domain.AccountConnection{
						AccountID:   "123456789",
						CountryCode: 280,
						BankID:      "1000000",
					},
				),
				UserID:          element.NewIdentification("user"),
				AccountCurrency: element.NewCurrency("EUR"),
				Name1:           element.NewAlphaNumeric("name", 27),
			},
			expectedAccount: domain.AccountInformation{
				AccountConnection: domain.AccountConnection{
					AccountID:   "123456789",
					CountryCode: 280,
					BankID:      "1000000",
				},
				UserID:   "user",
				Currency: "EUR",
				Name1:    "name",
			},
		},
		{
			desc: "v5 all mandatory fields set",
			segment: &AccountInformationV5{
				UserID: element.NewIdentification("user"),
				Name1:  element.NewAlphaNumeric("name", 27),
			},
			expectedAccount: domain.AccountInformation{
				UserID: "user",
				Name1:  "name",
			},
		},
		{
			desc: "v6 all mandatory fields set",
			segment: &AccountInformationV6{
				UserID: element.NewIdentification("user"),
				Name1:  element.NewAlphaNumeric("name", 27),
			},
			expectedAccount: domain.AccountInformation{
				UserID: "user",
				Name1:  "name",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			account := tt.segment.Account()

			if !reflect.DeepEqual(tt.expectedAccount, account) {
				t.Errorf("Expected accout to equal\n%v\n\tgot\n%v\n", tt.expectedAccount, account)
			}
		})
	}
}
