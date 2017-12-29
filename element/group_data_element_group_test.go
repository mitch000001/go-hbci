package element

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestAccountConnectionUnmarshalHBCI(t *testing.T) {
	test := "abc:subacc:280:12345678"

	acc := &AccountConnectionDataElement{}

	err := acc.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := NewAccountConnection(
		domain.AccountConnection{
			AccountID:                 "abc",
			SubAccountCharacteristics: "subacc",
			CountryCode:               280,
			BankID:                    "12345678",
		},
	).String()
	actual := acc.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
