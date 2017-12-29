package swift

import (
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/mitch000001/go-hbci/domain"
)

func TestAccountTagUnmarshal(t *testing.T) {
	// normal german BankID
	test := ":25:12345678/100000000"

	expected := &AccountTag{
		Tag:       ":25:",
		BankID:    "12345678",
		AccountID: "100000000",
	}

	tag := &AccountTag{}

	err := tag.Unmarshal([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expected, tag) {
		t.Logf("Expected tag to equal\n%#v\n\tgot\n%#v\n", expected, tag)
		t.Fail()
	}
}

func TestTransactionTagUnmarshal(t *testing.T) {
	tests := []struct {
		description    string
		marshaledValue string
		expectedTag    *TransactionTag
	}{
		{
			"All attributes set, booking in next month",
			":61:1511301202DR4,52N024NONREF//ABC\r\n/DEF",
			&TransactionTag{
				Tag:                   ":61:",
				ValutaDate:            domain.ShortDate{Time: domain.Date(2015, time.November, 30, time.Local).Truncate(24 * time.Hour)},
				BookingDate:           domain.ShortDate{Time: domain.Date(2015, time.December, 2, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator:  "D",
				CurrencyKind:          "R",
				Amount:                4.52,
				BookingKey:            "024",
				Reference:             "NONREF",
				BankReference:         "ABC",
				AdditionalInformation: "DEF",
			},
		},
		{
			"All attributes set, booking in new year",
			":61:1512300102DR4,52N024NONREF//ABC\r\n/DEF",
			&TransactionTag{
				Tag:                   ":61:",
				ValutaDate:            domain.ShortDate{Time: domain.Date(2015, time.December, 30, time.Local).Truncate(24 * time.Hour)},
				BookingDate:           domain.ShortDate{Time: domain.Date(2016, time.January, 2, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator:  "D",
				CurrencyKind:          "R",
				Amount:                4.52,
				BookingKey:            "024",
				Reference:             "NONREF",
				BankReference:         "ABC",
				AdditionalInformation: "DEF",
			},
		},
		{
			"All attributes set",
			":61:1508010803DR4,52N024NONREF//ABC\r\n/DEF",
			&TransactionTag{
				Tag:                   ":61:",
				ValutaDate:            domain.ShortDate{Time: domain.Date(2015, 8, 1, time.Local).Truncate(24 * time.Hour)},
				BookingDate:           domain.ShortDate{Time: domain.Date(2015, 8, 3, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator:  "D",
				CurrencyKind:          "R",
				Amount:                4.52,
				BookingKey:            "024",
				Reference:             "NONREF",
				BankReference:         "ABC",
				AdditionalInformation: "DEF",
			},
		},
		{
			"All attributes set except 'AdditionalInformation'",
			":61:1508010803DR4,52N024NONREF//ABC",
			&TransactionTag{
				Tag:                  ":61:",
				ValutaDate:           domain.ShortDate{Time: domain.Date(2015, 8, 1, time.Local).Truncate(24 * time.Hour)},
				BookingDate:          domain.ShortDate{Time: domain.Date(2015, 8, 3, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator: "D",
				CurrencyKind:         "R",
				Amount:               4.52,
				BookingKey:           "024",
				Reference:            "NONREF",
				BankReference:        "ABC",
			},
		},
		{
			"All attributes set except 'BankReference'",
			":61:1508010803DR4,52N024NONREF\r\n/DEF",
			&TransactionTag{
				Tag:                   ":61:",
				ValutaDate:            domain.ShortDate{Time: domain.Date(2015, 8, 1, time.Local).Truncate(24 * time.Hour)},
				BookingDate:           domain.ShortDate{Time: domain.Date(2015, 8, 3, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator:  "D",
				CurrencyKind:          "R",
				Amount:                4.52,
				BookingKey:            "024",
				Reference:             "NONREF",
				AdditionalInformation: "DEF",
			},
		},
		{
			"All attributes set except 'AdditionalInformation' and 'BankReference'",
			":61:1508010803DR4,52N024NONREF",
			&TransactionTag{
				Tag:                  ":61:",
				ValutaDate:           domain.ShortDate{Time: domain.Date(2015, 8, 1, time.Local).Truncate(24 * time.Hour)},
				BookingDate:          domain.ShortDate{Time: domain.Date(2015, 8, 3, time.Local).Truncate(24 * time.Hour)},
				DebitCreditIndicator: "D",
				CurrencyKind:         "R",
				Amount:               4.52,
				BookingKey:           "024",
				Reference:            "NONREF",
			},
		},
	}

	for _, test := range tests {
		tag := &TransactionTag{}

		err := tag.Unmarshal([]byte(test.marshaledValue))

		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		if !reflect.DeepEqual(test.expectedTag, tag) {
			pretty.Ldiff(t, test.expectedTag, tag)
			t.Logf("Expected tag to equal\n%#v\n\tgot\n%#v\n", test.expectedTag, tag)
			t.Fail()
		}
	}
}
