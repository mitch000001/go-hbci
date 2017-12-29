package segment

import (
	"testing"

	"github.com/mitch000001/go-hbci/domain"
)

func TestMessageHeaderSegmentUnmarshalHBCI(t *testing.T) {
	test := "HNHBK:1:3+000000000123+220+dialogID+3+'"

	header := &MessageHeaderSegment{}

	err := header.UnmarshalHBCI([]byte(test))
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := NewMessageHeaderSegment(123, 220, "dialogID", 3).String()
	actual := header.String()

	if expected != actual {
		t.Logf("Expected message header to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}

	// referencing header
	test = "HNHBK:1:3+000000000123+220+dialogID+3+abcde:1'"

	header = &MessageHeaderSegment{}

	err = header.UnmarshalHBCI([]byte(test))
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected = NewReferencingMessageHeaderSegment(
		123, 220, "dialogID", 3, domain.ReferencingMessage{DialogID: "abcde", MessageNumber: 1},
	).String()
	actual = header.String()

	if expected != actual {
		t.Logf("Expected message header to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
