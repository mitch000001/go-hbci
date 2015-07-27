package segment

import (
	"testing"

	"github.com/mitch000001/go-hbci/element"
)

func TestBankAnnouncementSegmentUnmarshalHBCI(t *testing.T) {
	test := "HIKIM:10:2+ec-Karte+Ihre neue ec-Karte liegt zur Abholung bereit.'"

	segment := &BankAnnouncementSegment{}

	expectedSegment := &BankAnnouncementSegment{
		Subject: element.NewAlphaNumeric("ec-Karte", 35),
		Body:    element.NewText("Ihre neue ec-Karte liegt zur Abholung bereit.", 2048),
	}
	expectedSegment.Segment = NewBasicSegment(10, expectedSegment)

	expected := expectedSegment.String()

	err := segment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	actual := segment.String()

	if expected != actual {
		t.Logf("Expected segment to equal\n%q\n\tgot:\n%q\n", expected, actual)
		t.Fail()
	}
}
