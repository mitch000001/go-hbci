package segment

import "testing"

func TestCommonBankParameterSegmentUnmarshalHBCI(t *testing.T) {
	t.Skip("TODO")
	test := "HIBPA:2:2:+12+280:10000000+Bank Name+3+1+201:210:220+0'"

	bankSegment := &CommonBankParameterSegment{}

	err := bankSegment.UnmarshalHBCI([]byte(test))

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}
}
