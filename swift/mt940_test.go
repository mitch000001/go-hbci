package swift

import (
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-hbci/domain"
)

func TestTransactionTagUnmarshal(t *testing.T) {
	test := ":61:1508010803DR4,52N024NONREF"

	expectedTransaction := &TransactionTag{
		Tag:                  ":61:",
		Date:                 domain.ShortDate{domain.Date(2015, 8, 1, time.Local).Truncate(24 * time.Hour)},
		BookingDate:          domain.ShortDate{domain.Date(2015, 8, 3, time.Local).Truncate(24 * time.Hour)},
		DebitCreditIndicator: "D",
		CurrencyKind:         "R",
		Amount:               4.52,
		BookingKey:           "024",
		Reference:            "NONREF",
	}

	tag := &TransactionTag{}

	err := tag.Unmarshal([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedTransaction, tag) {
		t.Logf("Expected transaction to equal\n%#v\n\tgot\n%#v\n", expectedTransaction, tag)
		t.Fail()
	}
}
