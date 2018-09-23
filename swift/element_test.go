package swift

import (
	"reflect"
	"testing"
)

func TestCustomFieldTagUnmarshal(t *testing.T) {
	test := ":86:123?00ABC?10xyz?20ah\r\nh?21hjj?301000?3156?32Max?33Muster?3499?60uu?61z\r\n4"

	expectedTag := &CustomFieldTag{
		Tag:                ":86:",
		TransactionID:      123,
		BookingText:        "ABC",
		PrimanotenNumber:   "xyz",
		Purpose:            []string{"ahh", "hjj"},
		BankID:             "1000",
		AccountID:          "56",
		Name:               "Max Muster",
		MessageKeyAddition: 99,
		Purpose2:           []string{"uu", "z4"},
	}

	tag := &CustomFieldTag{}

	err := tag.Unmarshal([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedTag, tag) {
		t.Logf("Expected tag to equal\n%#v\n\tgot\n%#v\n", expectedTag, tag)
		t.Fail()
	}
}
