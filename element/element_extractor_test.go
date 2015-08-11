package element

import (
	"reflect"
	"testing"
)

func TestElementExtractorExtract(t *testing.T) {
	type testData struct {
		in  string
		out []string
		err error
	}
	tests := []testData{
		{
			"abcde:123:012'",
			[]string{
				"abcde",
				"123",
				"012",
			},
			nil,
		},
		{
			"abcde:123:012+",
			[]string{
				"abcde",
				"123",
				"012",
			},
			nil,
		},
		{
			"abcde:123:012+de+1'",
			[]string{
				"abcde",
				"123",
				"012",
				"de",
				"1",
			},
			nil,
		},
		{
			"de+1+abcde:123:012'",
			[]string{
				"de",
				"1",
				"abcde",
				"123",
				"012",
			},
			nil,
		},
		{
			"abcde:123:012",
			[]string{
				"abcde",
				"123",
				"012",
			},
			nil,
		},
	}

	for _, test := range tests {
		extractor := NewElementExtractor([]byte(test.in))

		extracted, err := extractor.Extract()

		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		if extracted == nil {
			t.Logf("Expected result not to be nil\n")
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
