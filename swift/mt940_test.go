package swift

import (
	"reflect"
	"strconv"
	"strings"
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

		err := tag.Unmarshal([]byte(test.marshaledValue), 2015)

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

func TestTransactionTagOrder(t *testing.T) {
	testdata := "\r\n:20:HBCIKTOLST"
	for i := 0; i < 10; i++ {
		testdata += "\r\n:25:12345678/1234123456" +
			"\r\n:28C:0" +
			"\r\n:60F:C181105EUR1234,56" +
			"\r\n:61:1811051105DR50,NMSCNONREF" +
			"\r\n/OCMT/EUR50,//CHGS/   0,/" +
			"\r\n:86:177?00SB-SEPA-Ueberweisung?20" + strconv.Itoa(i+10) + "                                                                                                                                                 ?30?31?32Max Meier                  ?33                           ?34000" +
			"\r\n:62F:C190125EUR1234,56"

	}
	testdata += "\r\n-"
	mt := &MT940{}
	mt.Unmarshal([]byte(testdata))
	for i, tr := range mt.Transactions {
		if strings.TrimSpace(tr.Description.Purpose[0]) != strconv.Itoa(i+10) {
			t.Logf("Purpose at index %d should be %d but is %s", i, i+10, tr.Description.Purpose[0])
			t.Fail()
		}

	}

}

func TestTransactionListWithUnvalidData(t *testing.T) {
	testdata := "\r\n:20:HBCIKTOLST"
	testdata += "\r\n:25:12345678/1234123456" +
		"\r\n:28C:0" +
		"\r\n:60F:C181105EUR1234,56" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?20                                                                                                                                                   ?30?31?32Max Meier                  ?33                           ?34000" +
		"\r\n:62F:C190125EUR1234,56"

	testdata += "\r\n-"
	mt := &MT940{}
	error := mt.Unmarshal([]byte(testdata))
	if error == nil {
		t.Log("Error expected because of CustomTag without TransactionTag")
		t.Fail()
	}

}

func TestTransactionWithRedefineCustomDataTag(t *testing.T) {
	testdata := "\r\n:20:HBCIKTOLST"
	testdata += "\r\n:25:12345678/1234123456" +
		"\r\n:28C:0" +
		"\r\n:60F:C181105EUR1234,56" +

		"\r\n:61:1811051105DR50,NMSCNONREF" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?20                                                                                                                                                   ?30?31?32Max Meier                  ?33                           ?34000" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?20                                                                                                                                                   ?30?31?32Max Meier                  ?33                           ?34000" +
		"\r\n:62F:C190125EUR1234,56"

	testdata += "\r\n-"
	mt := &MT940{}
	error := mt.Unmarshal([]byte(testdata))
	if error == nil {
		t.Log("Error expected because of more than CustomTag after TransactionTag")
		t.Fail()
	}
}

func TestBookingDateBug(t *testing.T) {
	testdata := "\r\n:20:HBCIKTOLST"
	// Valuta Date 2019-01-01 Booking Date 2018-12-28
	testdata += "\r\n:25:12345678/1234123456" +
		"\r\n:28C:0" +
		"\r\n:60F:C181228EUR1234,56" +
		"\r\n:61:1901011228DR50,NMSCNONREF" +
		"\r\n/OCMT/EUR50,//CHGS/   0,/" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?202018-12-28                                                                                                                                           ?30?31?32Max Maier                  ?33                           ?34000" +
		"\r\n:62F:C190125EUR1234,56"
	// Valuta Date 2019-01-01 Booking Date 2019-01-01
	testdata += "\r\n:25:12345678/1234123456" +
		"\r\n:28C:0" +
		"\r\n:60F:C181228EUR1234,56" +
		"\r\n:61:1901010101DR50,NMSCNONREF" +
		"\r\n/OCMT/EUR50,//CHGS/   0,/" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?202019-01-01                                                                                                                                           ?30?31?32Max Meier                  ?33                           ?34000" +
		"\r\n:62F:C190125EUR1234,56"
	// Valuta Date 2018-12-28 Booking Date 2019-01-01 (not sure if that can happen ? )
	testdata += "\r\n:25:12345678/1234123456" +
		"\r\n:28C:0" +
		"\r\n:60F:C181228EUR1234,56" +
		"\r\n:61:1812280101DR50,NMSCNONREF" +
		"\r\n/OCMT/EUR50,//CHGS/   0,/" +
		"\r\n:86:177?00SB-SEPA-Ueberweisung?202019-01-01                                                                                                                                           ?30?31?32Max Meier                  ?33                           ?34000" +
		"\r\n:62F:C190125EUR1234,56"
	testdata += "\r\n-"
	mt := &MT940{}
	mt.Unmarshal([]byte(testdata))
	for _, tr := range mt.Transactions {
		if strings.TrimSpace(tr.Description.Purpose[0]) != tr.Transaction.BookingDate.String() {
			t.Logf("Booking date should be %s but is %s", strings.TrimSpace(tr.Description.Purpose[0]), tr.Transaction.BookingDate.String())
			t.Fail()
		}

	}

}
