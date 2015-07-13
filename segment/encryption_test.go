package segment

import "testing"

func TestEncryptedDataSegmentUnmarshalHBCI(t *testing.T) {
	test := "HNVSD:999:1+@5@abcde"

	enc := &EncryptedDataSegment{}

	err := enc.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := NewEncryptedDataSegment([]byte("abcde")).String()
	actual := enc.String()

	if expected != actual {
		t.Logf("Expected segment to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}
}
