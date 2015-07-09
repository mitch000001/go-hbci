package hbci

import (
	"reflect"
	"testing"
)

func TestSegmentExtratorExtract(t *testing.T) {
	type testData struct {
		in  string
		out []string
		err error
	}

	tests := []testData{
		{
			"HNHBK:1:3+abc'HNDGC:2:3+def'",
			[]string{
				"HNHBK:1:3+abc",
				"HNDGC:2:3+def",
			},
			nil,
		},
		{
			"HNHBK:1:3+abc'HNDGC:2:3+de?'f'",
			[]string{
				"HNHBK:1:3+abc",
				"HNDGC:2:3+de?'f",
			},
			nil,
		},
	}

	for _, test := range tests {
		extractor := NewSegmentExtractor([]byte(test.in))

		extracted, err := extractor.Extract()

		if test.err != nil && err == nil {
			t.Logf("Expected error, got nil\n")
			t.Fail()
		}

		if test.err == nil && err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		actual := make([]string, len(extracted))
		for i, b := range extracted {
			actual[i] = string(b)
		}

		if !reflect.DeepEqual(test.out, actual) {
			t.Logf("Extract: \n%q\n", extracted)
			t.Logf("Expected result to equal\n%q\n\tgot\n%q\n", test.out, actual)
			t.Fail()
		}
	}
}

func TestSegmentExtratorFindSegment(t *testing.T) {
	test := "HNHBK:1:3+abc'HNDGC:2:3+def'"

	extractor := NewSegmentExtractor([]byte(test))
	extractor.Extract()

	segment := extractor.FindSegment("HNHBK")

	if segment == nil {
		t.Logf("Expected segment not to be nil")
		t.Fail()
	}

	expected := []byte("HNHBK:1:3+abc")

	if !reflect.DeepEqual(expected, segment) {
		t.Logf("Expected segment to equal\n%q\n\tgot\n%q\n", expected, segment)
		t.Fail()
	}
}
