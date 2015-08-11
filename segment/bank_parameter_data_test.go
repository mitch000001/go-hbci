package segment

import (
	"reflect"
	"testing"

	"github.com/mitch000001/go-hbci/domain"
	"github.com/mitch000001/go-hbci/element"
)

func TestCommonBankParameterSegmentUnmarshalHBCI(t *testing.T) {
	test := "HIBPA:2:2:+12+280:10000000+Bank Name+3+1+201:210:220+0'"

	bankSegment := &CommonBankParameterSegment{}

	err := bankSegment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}
}

func TestBusinessTransactionParamsSegment(t *testing.T) {
	test := "DIDFBS:21:1:4+1+1+1'"

	segment := &BusinessTransactionParamsSegment{}

	expectedSegment := &BusinessTransactionParamsSegment{
		ID:            "DIDFBS",
		Version:       1,
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
