package segment

import "testing"

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
}
