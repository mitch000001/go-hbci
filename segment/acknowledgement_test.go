package segment

import "testing"

func TestMessageAcknowledgementUnmarshalHBCI(t *testing.T) {
	test := "HIRMG:1:2:+0010:1:Nachricht entgegengenommen:+0010:1:Nachricht entgegengenommen:"

	messageAcknowledgement := &MessageAcknowledgement{}

	err := messageAcknowledgement.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	marshaled := messageAcknowledgement.String()
	marshaled = marshaled[:len(marshaled)-1]

	if marshaled != test {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", test, marshaled)
		t.Fail()
	}
}

func TestSegmentAcknowledgementUnmarshalHBCI(t *testing.T) {
	test := "HIRMS:1:2:+0010:1:Nachricht entgegengenommen:+0010:1:Nachricht entgegengenommen:"

	segmentAcknowledgement := &SegmentAcknowledgement{}

	err := segmentAcknowledgement.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	marshaled := segmentAcknowledgement.String()
	marshaled = marshaled[:len(marshaled)-1]

	if marshaled != test {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", test, marshaled)
		t.Fail()
	}
}
