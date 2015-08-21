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

	v2 := &CommonUserParameterDataV2{
		UserID:     element.NewIdentification("12345"),
		UPDVersion: element.NewNumber(4, 3),
		UPDUsage:   element.NewNumber(0, 1),
	}
	v2.Segment = NewReferencingBasicSegment(5, 7, v2)
	expected := &CommonUserParameterDataSegment{v2}

	expectedString := expected.String()
	actualString := segment.String()

	if expectedString != actualString {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expectedString, actualString)
		t.Fail()
	}
}
