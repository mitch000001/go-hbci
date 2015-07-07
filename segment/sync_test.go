package segment

import "testing"

func TestSynchronisationSegmentUnmarshalHBCI(t *testing.T) {
	t.Skip()
	test := "HISYN:169:3:5+4534131272717070"

	syncMessage := &SynchronisationSegment{}

	err := syncMessage.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

}
