package segment

import (
	"testing"

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
