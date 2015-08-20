package segment

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestBusinessTransactionParamsSegment(t *testing.T) {
	test := "DIDFBS:21:1:4+1+1+1'"

	segment := &BusinessTransactionParamsSegment{}

	expectedSegment := &BusinessTransactionParamsSegment{
		id:            "DIDFBS",
		version:       1,
		MaxJobs:       element.NewNumber(1, 1),
		MinSignatures: element.NewNumber(1, 1),
	}
	expectedSegment.Segment = NewReferencingBasicSegment(21, 4, expectedSegment)
	expected := expectedSegment.String()

	err := segment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := segment.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}

func TestPinTanBusinessTransactionParamsSegmentUnmarshalHBCI(t *testing.T) {
	test := "DIPINS:3:1:4+1+1+HKSAL:N:HKUEB:J'"

	segment := &PinTanBusinessTransactionParamsSegment{}

	pinTanBusinessTransactions := []domain.PinTanBusinessTransaction{
		domain.PinTanBusinessTransaction{
			SegmentID: "HKSAL",
			NeedsTan:  false,
		},
		domain.PinTanBusinessTransaction{
			SegmentID: "HKUEB",
			NeedsTan:  true,
		},
	}
	pinTanDataElement := element.NewPinTanBusinessTransactionParameters(pinTanBusinessTransactions)
	expectedSegment := &PinTanBusinessTransactionParamsSegment{
		BusinessTransactionParamsSegment: &BusinessTransactionParamsSegment{
			MaxJobs:       element.NewNumber(1, 1),
			MinSignatures: element.NewNumber(1, 1),
			Params:        pinTanDataElement,
		},
	}
	expectedSegment.Segment = NewReferencingBasicSegment(3, 4, expectedSegment)
	expected := expectedSegment.String()

	err := segment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := segment.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}

func TestPinTanBusinessTransactionParamsSegmentPinTanBusinessTransactions(t *testing.T) {
	pinTanBusinessTransactions := []domain.PinTanBusinessTransaction{
		domain.PinTanBusinessTransaction{
			SegmentID: "HKSAL",
			NeedsTan:  false,
		},
		domain.PinTanBusinessTransaction{
			SegmentID: "HKUEB",
			NeedsTan:  true,
		},
	}
	pinTanDataElement := element.NewPinTanBusinessTransactionParameters(pinTanBusinessTransactions)
	segment := &PinTanBusinessTransactionParamsSegment{
		BusinessTransactionParamsSegment: &BusinessTransactionParamsSegment{
			MaxJobs:       element.NewNumber(1, 1),
			MinSignatures: element.NewNumber(1, 1),
			Params:        pinTanDataElement,
		},
	}
	segment.Segment = NewReferencingBasicSegment(3, 4, segment)

	expectedTransactions := []domain.PinTanBusinessTransaction{
		domain.PinTanBusinessTransaction{"HKSAL", false},
		domain.PinTanBusinessTransaction{"HKUEB", true},
	}

	pinTanTransactions := segment.PinTanBusinessTransactions()

	if !reflect.DeepEqual(expectedTransactions, pinTanTransactions) {
		t.Logf("Expected pinTanBusinessTransactions to return\n%+#v\n\tgot\n%+#v\n", expectedTransactions, pinTanTransactions)
		t.Fail()
	}
}
