package segment

import "testing"

func TestSynchronisationResponseSegmentUnmarshalHBCI(t *testing.T) {
	test := "HISYN:169:3:5+4534131272717070+12+144115'"

	syncSegment := &SynchronisationResponseSegment{}

	err := syncSegment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	expected := test
	actual := syncSegment.String()

	if expected != actual {
		t.Logf("Expected unmarshaled value to equal\n%q\n\tgot\n%q\n", expected, actual)
		t.Fail()
	}

}
