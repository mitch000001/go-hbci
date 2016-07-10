package spec

import "testing"

func TestGenerate(t *testing.T) {
	groups := []DataElementSpec{}
	segments := []SegmentSpec{}

	expected := ""

	actual, err := Generate(groups, segments)

	if err != nil {
		t.Logf("Expected no error, got: %v\n", err)
		t.Fail()
	}

	if expected != string(actual) {
		t.Logf("Expected result to equal\n%q\n\tgot:\n%q\n", expected, string(actual))
		t.Fail()
	}
}
