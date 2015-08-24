package message

import (
	"testing"

	"github.com/mitch000001/go-hbci/element"
	"github.com/mitch000001/go-hbci/segment"
)

func TestUnmarshalerUnmarshalSegment(t *testing.T) {
	test := "HNHBK:1:3+000000000273+220+abcde+1+'HISYN:1:3+abcde++'"

	unmarshaler := NewUnmarshaler([]byte(test))

	seg, err := unmarshaler.UnmarshalSegment("HNHBK", 3)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := segment.NewMessageHeaderSegment(273, 220, "abcde", 1).String()

	if seg != nil {
		actual := seg.String()

		if expected != actual {
			t.Logf("Expected segment to equal\n%q\n\tgot\n%q\n", expected, actual)
			t.Fail()
		}
	} else {
		t.Logf("Expected segment not to be nil\n")
		t.Fail()
	}

	// Test another segment
	seg, err = unmarshaler.UnmarshalSegment("HISYN", 3)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	syncSegment := &segment.SynchronisationResponseSegment{ClientSystemID: element.NewIdentification("abcde")}
	syncSegment.Segment = segment.NewBasicSegment(1, syncSegment)
	expected = syncSegment.String()

	if seg != nil {
		actual := seg.String()

		if expected != actual {
			t.Logf("Expected segment to equal\n%q\n\tgot\n%q\n", expected, actual)
			t.Fail()
		}
	} else {
		t.Logf("Expected segment not to be nil\n")
		t.Fail()
	}

	// Test unknown segment
	test = "HXXXX:1:3+abcde++'"

	unmarshaler = NewUnmarshaler([]byte(test))

	seg, err = unmarshaler.UnmarshalSegment("HXXXX", 3)

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		errMessage := err.Error()
		expectedMessage := "Segment not in index: \"HXXXX:3\""
		if expectedMessage != errMessage {
			t.Logf("Expected message to equal %s, got %s\n", expectedMessage, errMessage)
			t.Fail()
		}
	}

	// Test known but not included segment
	test = "HXXXX:1:3+abcde++'"

	unmarshaler = NewUnmarshaler([]byte(test))

	seg, err = unmarshaler.UnmarshalSegment("HNHBK", 3)

	if err == nil {
		t.Logf("Expected error, got nil\n")
		t.Fail()
	} else {
		errMessage := err.Error()
		expectedMessage := "Segment not found in message: \"HNHBK\""
		if expectedMessage != errMessage {
			t.Logf("Expected message to equal %q, got %q\n", expectedMessage, errMessage)
			t.Fail()
		}
	}
}
